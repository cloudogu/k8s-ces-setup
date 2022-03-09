# Set these to the desired values
ARTIFACT_ID=k8s-ces-setup
VERSION=0.0.0
GOTAG?=1.17.7
MAKEFILES_VERSION=4.8.0

# Image URL to use all building/pushing image targets
IMAGE ?= cloudogu/${ARTIFACT_ID}:${VERSION}

K8S_RESOURCE_DIR=${WORKDIR}/k8s
K8S_SETUP_CONFIG_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-ces-setup-config.yaml
K8S_SETUP_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-ces-setup.yaml
K8S_SETUP_DEV_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-ces-setup_dev.yaml
K8S_CLUSTER_ROOT=<your/path/to/k3ces>

.DEFAULT_GOAL:=help

include build/make/variables.mk

ADDITIONAL_CLEAN=clean-vendor
PRE_COMPILE=vet

include build/make/self-update.mk
include build/make/info.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ EcoSystem

.PHONY: build
build: docker-build image-import k8s-apply ## Builds a new version of the setup and deploys it into the K8s-EcoSystem.

##@ Development (without go container)

.PHONY: vet
vet: $(STATIC_ANALYSIS_DIR) ## Run go vet against code.
	@go vet ./... | tee ${STATIC_ANALYSIS_DIR}/report-govet.out

##@ Build

.PHONY: build-setup
build-setup: ## Builds the setup Go binary.
# pseudo target to support make help for compile target
	@make compile

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
docker-build: ## Builds the docker image of the k8s-ces-setup `registry.cloudogu.com/official/k8s-ces-setup:<version>>`.
	@echo "Building docker image of dogu..."
	docker build . -t ${IMAGE}

${K8S_CLUSTER_ROOT}/image.tar: # [not listed in help] Saves the `registry.cloudogu.com/official/nginx-ingress:dev` image into a file into the K8s root path to be available on all nodes.
	docker save ${IMAGE} -o ${K8S_CLUSTER_ROOT}/image.tar

.PHONY: image-import
image-import: ${K8S_CLUSTER_ROOT}/image.tar ## Imports the currently available image `registry.cloudogu.com/official/nginx-ingress:dev` into the K8s cluster for all nodes.
	@echo "Import docker image of dogu into all K8s nodes..."
	cd ${K8S_CLUSTER_ROOT} && vagrant ssh main -- -t "sudo k3s ctr images import /vagrant/image.tar"
	cd ${K8S_CLUSTER_ROOT} && vagrant ssh worker-0 -- -t "sudo k3s ctr images import /vagrant/image.tar"
	cd ${K8S_CLUSTER_ROOT} && vagrant ssh worker-1 -- -t "sudo k3s ctr images import /vagrant/image.tar"
	rm ${K8S_CLUSTER_ROOT}/image.tar

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

${K8S_SETUP_DEV_RESOURCE_YAML}: # [not listed in help] Templates the deployment yaml with the latest image.
	@yq e "(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.name == \"k8s-ces-setup\")).image=\"${IMAGE}\"" ${K8S_SETUP_RESOURCE_YAML} > ${K8S_SETUP_DEV_RESOURCE_YAML}

# Other targets

.PHONY: clean-vendor
clean-vendor:
	rm -rf vendor
