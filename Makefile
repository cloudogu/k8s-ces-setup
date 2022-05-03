# Set these to the desired values
ARTIFACT_ID=k8s-ces-setup
VERSION=0.3.0

GOTAG?=1.17.7
MAKEFILES_VERSION=5.1.0

# Image URL to use all building/pushing image targets
IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}

K8S_RESOURCE_DIR=${WORKDIR}/k8s
K8S_SETUP_CONFIG_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-ces-setup-config.yaml
K8S_SETUP_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-ces-setup.yaml

LOCAL_HTTP_SERVER_PORT=9876
LOCAL_HTTP_DIR=k8s/dev-resources

include build/make/variables.mk

# make sure to create a statically linked binary otherwise it may quit with
# "exec user process caused: no such file or directory"
GO_BUILD_FLAGS=-mod=vendor -a -tags netgo,osusergo $(LDFLAGS) -o $(BINARY)
# remove DWARF symbol table and strip other symbols to shave ~13 MB from binary
ADDITIONAL_LDFLAGS=-extldflags -static -w -s

include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
include build/make/k8s.mk

##@ EcoSystem

.PHONY: build
build: k8s-delete image-import k8s-apply ## Builds a new version of the setup and deploys it into the K8s-EcoSystem.

##@ Development (without go container)

.PHONY: serve-local-yaml
serve-local-yaml:
	@echo "Starting to server ${WORKDIR}/${LOCAL_HTTP_DIR}"
	@echo "Press Ctrl+C to exit"
	@echo "You need a routable IP address or DNS in order to address resources from inside the cluster"
	@cd ${WORKDIR}/${LOCAL_HTTP_DIR} && python3 -m http.server ${LOCAL_HTTP_SERVER_PORT}

##@ Build

.PHONY: build-setup
build-setup: ${SRC} compile ## Builds the setup Go binary.

.PHONY: run
run: vet ## Run a setup from your host.
	go run ./main.go

.PHONY: k8s-create-temporary-resource
k8s-create-temporary-resource:
	@cp $(K8S_SETUP_CONFIG_RESOURCE_YAML) $(K8S_RESOURCE_TEMP_YAML)
	@echo "---" >> $(K8S_RESOURCE_TEMP_YAML)
	@cat $(K8S_SETUP_RESOURCE_YAML) >> $(K8S_RESOURCE_TEMP_YAML)

##@ Release

.PHONY: setup-release
setup-release: ## Interactively starts the release workflow.
	@echo "Starting git flow release..."
	@build/make/release.sh setup