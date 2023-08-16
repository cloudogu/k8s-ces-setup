# Set these to the desired values
ARTIFACT_ID=k8s-ces-setup
VERSION=0.16.0

GOTAG?=1.20.4
MAKEFILES_VERSION=7.10.0

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## Image URL to use all building/pushing image targets
IMAGE_DEV=${K3CES_REGISTRY_URL_PREFIX}/${ARTIFACT_ID}:${VERSION}
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
LINT_VERSION?=v1.52.1

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
include build/make/mocks.mk
include build/make/release.mk

K8S_PRE_GENERATE_TARGETS=k8s-create-temporary-resource template-dev-only-image-pull-policy

#TODO test pw
HELM_DOGU_REGISTRY_ARGS=--set=dogu_registry_secret.username='${DOGU_REGISTRY_USERNAME}' --set=dogu_registry_secret.password='${DOGU_REGISTRY_PASSWORD}'
HELM_DOCKER_REGISTRY_ARGS=--set=docker_registry_secret.username='${DOCKER_REGISTRY_USERNAME}' --set=docker_registry_secret.password='${DOCKER_REGISTRY_PASSWORD}'
HELM_SETUP_JSON_ARGS=--set-file=setup_json="${WORKDIR}/k8s/dev-resources/setup.json"
ADDITIONAL_HELM_APPLY_ARGS=${HELM_DOGU_REGISTRY_ARGS} ${HELM_DOCKER_REGISTRY_ARGS} ${HELM_SETUP_JSON_ARGS}

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

##@ Development (with cluster)

.PHONY: k8s-clean
k8s-clean: ## Cleans all resources deployed by the setup
	@kubectl delete --all dogus --namespace=$(K8S_CURRENT_NAMESPACE) || true
	@kubectl delete all -l app.kubernetes.io/name=k8s-ces-setup --namespace=$(K8S_CURRENT_NAMESPACE) || true
	@kubectl delete all -l app.kubernetes.io/name=k8s-ces-setup-finisher --namespace=$(K8S_CURRENT_NAMESPACE) || true
	@kubectl delete all -l app.kubernetes.io/name=etcd --namespace=$(K8S_CURRENT_NAMESPACE) || true
	@kubectl delete all -l run=etcd-client --namespace=$(K8S_CURRENT_NAMESPACE) || true
	@kubectl delete all -l control-plane=controller-manager --namespace=$(K8S_CURRENT_NAMESPACE) || true
	@kubectl delete secret k8s-dogu-operator-dogu-registry || true
	@kubectl delete secret docker-registry k8s-dogu-operator-docker-registry || true
	@kubectl get clusterroles,clusterrolebindings | grep k8s-dogu-operator | sed 's| .*||g' | xargs kubectl delete - || true
	@kubectl get clusterroles,clusterrolebindings | grep k8s-service-discovery | sed 's| .*||g' | xargs kubectl delete - || true
	@kubectl create ns $(K8S_CURRENT_NAMESPACE) || true
	@kubectl ns $(K8S_CURRENT_NAMESPACE)
	@kubectl create secret generic k8s-dogu-operator-dogu-registry --from-literal=endpoint=${DOGU_REGISTRY_URL} --from-literal=username=${DOGU_REGISTRY_USERNAME} --from-literal=password=${DOGU_REGISTRY_PASSWORD}
	@kubectl create secret docker-registry k8s-dogu-operator-docker-registry --docker-server=${DOCKER_REGISTRY_URL} --docker-username=${DOCKER_REGISTRY_USERNAME} --docker-email="" --docker-password=${DOCKER_REGISTRY_PASSWORD}
	@make build

##@ Build

.PHONY: build-setup
build-setup: ${SRC} compile ## Builds the setup Go binary.

.PHONY: setup-etcd-port-forward
setup-etcd-port-forward:
	kubectl port-forward etcd-0 4001:2379 &

.PHONY: run
run: ## Run a setup from your host.
	go run ./main.go

.PHONY: copy-setup-resources
copy-setup-resources:
	@cp $(K8S_SETUP_RESOURCE_YAML) $(K8S_RESOURCE_TEMP_YAML)

.PHONY: k8s-create-temporary-resource
k8s-create-temporary-resource: $(K8S_RESOURCE_TEMP_FOLDER) copy-setup-resources template-dev-only-image-pull-policy
	@$(BINARY_YQ) -i e "(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.image == \"*$(ARTIFACT_ID)*\").image)=\"$(IMAGE)\"" $(K8S_RESOURCE_TEMP_YAML)

.PHONY: create-temporary-dev-resource
create-temporary-dev-resource: $(K8S_RESOURCE_TEMP_FOLDER) k8s-create-temporary-resource template-dev-only-image-pull-policy
	@echo "---" >> $(K8S_RESOURCE_TEMP_YAML)
	@kubectl create configmap k8s-ces-setup-json --from-file=k8s/dev-resources/setup.json --dry-run=client -o yaml >> $(K8S_RESOURCE_TEMP_YAML)
	@echo "---" >> $(K8S_RESOURCE_TEMP_YAML)
	@cat $(K8S_SETUP_CONFIG_RESOURCE_YAML) >> $(K8S_RESOURCE_TEMP_YAML)

.PHONY: template-dev-only-image-pull-policy
template-dev-only-image-pull-policy: $(BINARY_YQ)
	@if [[ ${STAGE}"X" == "development""X" ]]; \
		then echo "Setting pull policy to always for development stage!" && $(BINARY_YQ) -i e "(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.image == \"*$(ARTIFACT_ID)*\").imagePullPolicy)=\"Always\"" $(K8S_RESOURCE_TEMP_YAML); \
	fi
