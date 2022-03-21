# Set these to the desired values
ARTIFACT_ID=k8s-ces-setup
VERSION=0.1.0-dev

GOTAG?=1.17.7
MAKEFILES_VERSION=5.0.0

# Image URL to use all building/pushing image targets
IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}

K8S_RESOURCE_DIR=${WORKDIR}/k8s
K8S_SETUP_CONFIG_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-ces-setup-config.yaml
K8S_SETUP_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-ces-setup.yaml
K8S_SETUP_DEV_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-ces-setup_dev.yaml

LOCAL_HTTP_SERVER_PORT=9876
LOCAL_HTTP_DIR=k8s/dev-resources

include build/make/variables.mk

# make sure to create a statically linked binary otherwise it may quit with
# "exec user process caused: no such file or directory"
GO_BUILD_FLAGS=-mod=vendor -a -tags netgo,osusergo $(LDFLAGS) -o $(BINARY)
# remove DWARF symbol table and strip other symbols to shave ~13 MB from binary
ADDITIONAL_LDFLAGS=-extldflags -static -w -s

.DEFAULT_GOAL:=help

include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk

##@ EcoSystem

.PHONY: build
build: docker-build image-import k8s-apply ## Builds a new version of the setup and deploys it into the K8s-EcoSystem.

##@ Development (without go container)

.PHONY: vet
vet: ${STATIC_ANALYSIS_DIR}/report-govet.out ## Run go vet against code.

${STATIC_ANALYSIS_DIR}/report-govet.out: ${SRC} $(STATIC_ANALYSIS_DIR)
	@go vet ./... | tee $@

##@ Build

.PHONY: build-setup
build-setup: ${SRC} compile ## Builds the setup Go binary.

.PHONY: run
run: vet ## Run a controller from your host.
	go run ./main.go

##@ Release

.PHONY: setup-release
setup-release: ## Interactively starts the release workflow.
	@echo "Starting git flow release..."
	@build/make/release.sh setup

##@ Docker

.PHONY: docker-build
docker-build: ${SRC} compile ## Builds the docker image of the k8s-ces-setup `cloudogu/k8s-ces-setup:version`.
	@echo "Building docker image of dogu..."
	docker build . -t ${IMAGE}

${K8S_CLUSTER_ROOT}/image.tar: check-k8s-cluster-root-env-var
	# Saves the `cloudogu/k8s-ces-setup:version` image into a file into the K8s root path to be available on all nodes.
	docker save ${IMAGE} -o ${K8S_CLUSTER_ROOT}/image.tar

.PHONY: image-import
image-import: ${K8S_CLUSTER_ROOT}/image.tar
    # Imports the currently available image `cloudogu/k8s-ces-setup:version` into the K8s cluster for all nodes.
	@echo "Import docker image of dogu into all K8s nodes..."
	@cd ${K8S_CLUSTER_ROOT} && \
		for node in $$(vagrant status --machine-readable | grep "state,running" | awk -F',' '{print $$2}'); \
		do  \
			echo "...$${node}"; \
			vagrant ssh $${node} -- -t "sudo k3s ctr images import /vagrant/image.tar"; \
		done;
	@echo "Done."
	rm ${K8S_CLUSTER_ROOT}/image.tar

.PHONY: check-k8s-cluster-root-env-var
check-k8s-cluster-root-env-var:
	@echo "Checking if env var K8S_CLUSTER_ROOT is set..."
	@bash -c export -p | grep K8S_CLUSTER_ROOT
	@echo "Done."

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: k8s-apply
k8s-apply: ${K8S_SETUP_DEV_RESOURCE_YAML} ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	@echo "Apply all k8s-ces-setup resources to the K8s-EcoSystem..."
	kubectl apply -f ${K8S_SETUP_CONFIG_RESOURCE_YAML}
	kubectl apply -f ${K8S_SETUP_DEV_RESOURCE_YAML}
	@rm ${K8S_SETUP_DEV_RESOURCE_YAML}

.PHONY: k8s-delete
k8s-delete: ${K8S_SETUP_DEV_RESOURCE_YAML} ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	@echo "Delete all k8s-ces-setup resources from the K8s-EcoSystem..."
	kubectl delete --ignore-not-found=true -f ${K8S_SETUP_CONFIG_RESOURCE_YAML}
	kubectl delete --ignore-not-found=true -f ${K8S_SETUP_DEV_RESOURCE_YAML}
	@rm ${K8S_SETUP_DEV_RESOURCE_YAML}

${K8S_SETUP_DEV_RESOURCE_YAML}:
	# Templates the deployment yaml with the latest image.
	@yq "(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.name == \"k8s-ces-setup\")).image=\"${IMAGE}\"" ${K8S_SETUP_RESOURCE_YAML} > ${K8S_SETUP_DEV_RESOURCE_YAML}

.PHONY: serve-local-yaml
serve-local-yaml:
	@echo "Starting to server ${WORKDIR}/${LOCAL_HTTP_DIR}"
	@echo "Press Ctrl+C to exit"
	@echo "You need a routable IP address or DNS in order to address resources from inside the cluster"
	@cd ${WORKDIR}/${LOCAL_HTTP_DIR} && python3 -m http.server ${LOCAL_HTTP_SERVER_PORT}