#!groovy

@Library(['github.com/cloudogu/dogu-build-lib@v1.4.1', 'github.com/cloudogu/ces-build-lib@v1.49.0'])
import com.cloudogu.ces.cesbuildlib.*
import com.cloudogu.ces.dogubuildlib.*

// Creating necessary git objects, object cannot be named 'git' as this conflicts with the method named 'git' from the library
gitWrapper = new Git(this, "cesmarvin")
gitWrapper.committerName = 'cesmarvin'
gitWrapper.committerEmail = 'cesmarvin@cloudogu.com'
gitflow = new GitFlow(this, gitWrapper)
github = new GitHub(this, gitWrapper)
changelog = new Changelog(this)
Docker docker = new Docker(this)
gpg = new Gpg(this, docker)

// Create necessary k3d object
k3d = new K3d(this, env.WORKSPACE, env.PATH)

// Configuration of repository
repositoryOwner = "cloudogu"
repositoryName = "k8s-ces-setup"
project = "github.com/${repositoryOwner}/${repositoryName}"

// Configuration of branches
productionReleaseBranch = "main"
developmentBranch = "develop"
currentBranch = "${env.BRANCH_NAME}"

node('docker') {
    timestamps {
        properties([
                // Keep only the last x builds to preserve space
                buildDiscarder(logRotator(numToKeepStr: '10')),
                // Don't run concurrent builds for a branch, because they use the same workspace directory
                disableConcurrentBuilds(),
        ])

        stage('Checkout') {
            gitWrapper.git branch: 'main', url: 'https://github.com/cloudogu/gitops-playground'
            dir('k8s-ces-setup') {
                checkout scm
                make 'clean'
            }
        }

        dir('k8s-ces-setup') {
            stage('Lint - Dockerfile') {
                lintDockerfile()
            }

            stage("Lint - k8s Resources") {
                stageLintK8SResources()
            }

            docker
                    .image('golang:1.17.7')
                    .mountJenkinsUser()
                    .inside("--volume ${PWD}:/go/src/${project} -w /go/src/${project}")
                            {
                                stage('Build') {
                                    make 'build'
                                }

                                stage('Unit Tests') {
                                    make 'unit-test'
                                }

                                stage("Review dog analysis") {
                                    stageStaticAnalysisReviewDog()
                                }
                            }

            stage('SonarQube') {
                stageStaticAnalysisSonarQube()
            }
        }

        try {
            stage('Set up k3d cluster') {
                k3d.startK3d()
            }

            stage('Install kubectl') {
                k3d.installKubectl()
            }

            stage('Build Image') {
                dir('k8s-ces-setup') {
                    make "docker-build"
                }
            }

            stage('Import Image') {
                String currentVersion = "dev"
                dir('k8s-ces-setup') {
                    currentVersion = getCurrentVersionFromMakefile()
                }
                sh "k3d image import ${repositoryOwner}/${repositoryName}:${currentVersion}"
            }

            stage('Deploy Setup') {
                k3d.kubectl("apply -f k8s-ces-setup/k8s/k8s-ces-setup.yaml")
            }

            dir('k8s-ces-setup') {
                stageAutomaticRelease()
            }
        } finally {
            stage('Remove k3d cluster') {
                k3d.deleteK3d()
            }
        }
    }
}

void gitWithCredentials(String command) {
    withCredentials([usernamePassword(credentialsId: 'cesmarvin', usernameVariable: 'GIT_AUTH_USR', passwordVariable: 'GIT_AUTH_PSW')]) {
        sh(
                script: "git -c credential.helper=\"!f() { echo username='\$GIT_AUTH_USR'; echo password='\$GIT_AUTH_PSW'; }; f\" " + command,
                returnStdout: true
        )
    }
}

void stageLintK8SResources() {
    String kubevalImage = "cytopia/kubeval:0.13"
    sh "printenv"
    String currentWorkspace = "${PWD}"
    sh "echo ${env.PWD}"
    sh "echo ${currentWorkspace}"
    docker
            .image(kubevalImage)
            .inside("-v ${currentWorkspace}/k8s:/data -t --entrypoint=")
                    {
                        sh "kubeval /data/k8s-ces-setup.yaml --ignore-missing-schemas"
                    }
}

void stageStaticAnalysisReviewDog() {
    def commitSha = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-gh', usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
        withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=cloudogu", "CI_REPO_NAME=${repositoryName}"]) {
            make 'static-analysis-ci'
        }
    }
}

void stageStaticAnalysisSonarQube() {
    def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
    withSonarQubeEnv {
        sh "git config 'remote.origin.fetch' '+refs/heads/*:refs/remotes/origin/*'"
        gitWithCredentials("fetch --all")

        if (currentBranch == productionReleaseBranch) {
            echo "This branch has been detected as the production branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (currentBranch == developmentBranch) {
            echo "This branch has been detected as the development branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (env.CHANGE_TARGET) {
            echo "This branch has been detected as a pull request."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.pullrequest.key=${env.CHANGE_ID} -Dsonar.pullrequest.branch=${env.CHANGE_BRANCH} -Dsonar.pullrequest.base=${developmentBranch}"
        } else if (currentBranch.startsWith("feature/")) {
            echo "This branch has been detected as a feature branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else {
            echo "This branch has been detected as a miscellaneous branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} "
        }
    }
    timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
        def qGate = waitForQualityGate()
        if (qGate.status != 'OK') {
            unstable("Pipeline unstable due to SonarQube quality gate failure")
        }
    }
}

void stageAutomaticRelease() {
    if (gitflow.isReleaseBranch()) {
        String releaseVersion = gitWrapper.getSimpleBranchName()

        stage('Build & Push Image') {
            def dockerImage = docker.build("cloudogu/${repositoryName}:${releaseVersion}")

            docker.withRegistry('https://registry.hub.docker.com/', 'dockerHubCredentials') {
                dockerImage.push("${releaseVersion}")
            }
        }

        stage('Finish Release') {
            gitflow.finishRelease(releaseVersion, productionReleaseBranch)
        }

        stage('Sign after Release') {
            gpg.createSignature()
        }

        stage('Add Github-Release') {
            releaseId = github.createReleaseWithChangelog(releaseVersion, changelog, productionReleaseBranch)
        }
    }
}

void make(String makeArgs) {
    sh "make ${makeArgs}"
}

String getCurrentVersionFromMakefile() {
    return sh(returnStdout: true, script: 'cat Makefile | grep -e "^VERSION=" | sed "s/VERSION=//g"').trim()
}