#!groovy
@Library('github.com/cloudogu/ces-build-lib@99c84661de69c44716139ba49747219b90be58a9')
import com.cloudogu.ces.cesbuildlib.*

// Creating necessary git objects, object cannot be named 'git' as this conflicts with the method named 'git' from the library
gitWrapper = new Git(this, "cesmarvin")
gitWrapper.committerName = 'cesmarvin'
gitWrapper.committerEmail = 'cesmarvin@cloudogu.com'
gitflow = new GitFlow(this, gitWrapper)
github = new GitHub(this, gitWrapper)
changelog = new Changelog(this)
Docker docker = new Docker(this)
goVersion = "1.20"

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
            checkout scm
            make 'dist-clean'
        }

        stage('Lint - Dockerfile') {
            lintDockerfile()
        }

        stage("Lint - k8s Resources") {
            stageLintK8SResources()
        }

        stage('Check Markdown Links') {
            Markdown markdown = new Markdown(this, "3.11.0")
            markdown.check()
        }

        docker
                .image("golang:${goVersion}")
                .mountJenkinsUser()
                .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}")
                        {
                            stage('Build') {
                                make 'compile'
                            }

                            stage('Unit Tests') {
                                make 'unit-test'
                                junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
                            }

                            stage("Review dog analysis") {
                                stageStaticAnalysisReviewDog()
                            }
                        }

        stage('SonarQube') {
            stageStaticAnalysisSonarQube()
        }

        def k3d = new K3d(this, "${WORKSPACE}", "${WORKSPACE}/k3d", env.PATH)

        try {
            stage('Set up k3d cluster') {
                k3d.startK3d()
            }

            def cessetupImageName
            stage('Build & Push Image') {
                def makefile = new Makefile(this)
                String setupVersion = makefile.getVersion()
                cessetupImageName = k3d.buildAndPushToLocalRegistry("cloudogu/${repositoryName}", setupVersion)
            }

            def sourceDeploymentYaml = "k8s/k8s-ces-setup.yaml"

            stage('Patch setup YAML to use local image') {
                docker.image('mikefarah/yq:4.22.1')
                        .mountJenkinsUser()
                        .inside("--volume ${WORKSPACE}:/workdir -w /workdir") {
                            sh "yq -i '(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.name == \"k8s-ces-setup\")).image=\"${cessetupImageName}\"' ${sourceDeploymentYaml}"
                            // avoid RBAC errors during installing the CRD because of empty ns vs default ns in the ClusterRoleBinding
                            sh "sed -i 's/{{ .Namespace }}/default/g' ${sourceDeploymentYaml}"
                        }
            }

            stage('Configure Setup') {
                k3d.assignExternalIP()
                def commitSha = getCurrentCommit()
                k3d.configureSetup(commitSha, [
                        dependencies: ["k8s/nginx-ingress"],
                        defaultDogu : ""
                ])
            }

            stage('Install and Trigger Setup (trigger warning: setup)') {
                k3d.kubectl("apply -f ${sourceDeploymentYaml}")
            }

            stage("wait for k8s-specific dogu (it has special needs)") {
                k3d.waitForDeploymentRollout("nginx-ingress", 300, 10)
            }

            stage('Restore development resources') {
                sh "git restore ${sourceDeploymentYaml}"
            }

            stageAutomaticRelease()
        } catch(Exception e) {
            k3d.collectAndArchiveLogs()
            throw e
        } finally {
            stage('Remove k3d cluster') {
                k3d.deleteK3d()
            }
        }
    }
}

void stageLintK8SResources() {
    String kubevalImage = "cytopia/kubeval:0.13"
    docker
            .image(kubevalImage)
            .inside("-v ${WORKSPACE}/k8s:/data -t --entrypoint=")
                    {
                        sh "kubeval /data/k8s-ces-setup.yaml --ignore-missing-schemas"
                    }
}

String getCurrentCommit() {
    return sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
}

void stageStaticAnalysisReviewDog() {
    def commitSha=getCurrentCommit()
    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-gh', usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
        withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=${repositoryOwner}", "CI_REPO_NAME=${repositoryName}"]) {
            make 'static-analysis'
        }
    }
}

void stageStaticAnalysisSonarQube() {
    def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
    withSonarQubeEnv {
        gitWrapper.fetch()

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
        Makefile makefile = new Makefile(this)
        String setupVersion = makefile.getVersion()

        stage('Build & Push Image') {
            def dockerImage = docker.build("cloudogu/${repositoryName}:${setupVersion}")

            docker.withRegistry('https://registry.hub.docker.com/', 'dockerHubCredentials') {
                dockerImage.push("${setupVersion}")
            }
        }

        stage('Finish Release') {
            gitflow.finishRelease(releaseVersion, productionReleaseBranch)
        }

        stage('Add Github-Release') {
            releaseId = github.createReleaseWithChangelog(releaseVersion, changelog, productionReleaseBranch)
        }

        stage('Regenerate resources for release') {
            make 'create-temporary-release-resource'
        }

        stage('Push to Registry') {
            GString targetSetupResourceYaml = "target/make/k8s/${repositoryName}_${setupVersion}.yaml"

            DoguRegistry registry = new DoguRegistry(this)
            registry.pushK8sYaml(targetSetupResourceYaml, repositoryName, "k8s", "${setupVersion}")
        }
    }
}

void make(String makeArgs) {
    sh "make ${makeArgs}"
}