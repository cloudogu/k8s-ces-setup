# Set these to the desired values
ARTIFACT_ID=k8s-ces-setup
VERSION=4.0.0

GOTAG?=1.24.3
MAKEFILES_VERSION=9.9.1

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

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
include build/make/mocks.mk
include build/make/release.mk

BINARY_HELM_ADDITIONAL_UPGR_ARGS=--set-file=setup_json="${WORKDIR}/k8s/dev-resources/setup.json" --values="${WORKDIR}/additionalValues.yaml"

K8S_COMPONENT_SOURCE_VALUES = ${HELM_SOURCE_DIR}/values.yaml
K8S_COMPONENT_TARGET_VALUES = ${HELM_TARGET_DIR}/values.yaml
HELM_PRE_APPLY_TARGETS=template-stage template-log-level template-image-pull-policy template-dogu-registry template-docker-registry template-helm-registry
HELM_PRE_GENERATE_TARGETS = helm-values-update-image-version
HELM_POST_GENERATE_TARGETS = helm-values-replace-image-repo
CHECK_VAR_TARGETS=check-all-vars
IMAGE_IMPORT_TARGET=image-import

include build/make/k8s-component.mk

##@ EcoSystem

.PHONY: build
build: helm-delete helm-apply ## Builds a new version of the setup and deploys it into the K8s-EcoSystem.

.PHONY: helm-values-update-image-version
helm-values-update-image-version: $(BINARY_YQ)
	@echo "Updating the image version in source values.yaml to ${VERSION}..."
	@$(BINARY_YQ) -i e ".setup.image.tag = \"${VERSION}\"" ${K8S_COMPONENT_SOURCE_VALUES}

.PHONY: helm-values-replace-image-repo
helm-values-replace-image-repo: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
      		echo "Setting dev image repo in target values.yaml!" ;\
    		$(BINARY_YQ) -i e ".setup.image.registry=\"$(shell echo '${IMAGE_DEV}' | sed 's/\([^\/]*\)\/\(.*\)/\1/')\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
    		$(BINARY_YQ) -i e ".setup.image.repository=\"$(shell echo '${IMAGE_DEV}' | sed 's/\([^\/]*\)\/\(.*\)/\2/')\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
    	fi

.PHONY: template-stage
template-stage: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
  		echo "Setting STAGE env in deployment to ${STAGE}!" ;\
		$(BINARY_YQ) -i e ".setup.manager.env.stage=\"${STAGE}\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
	fi

.PHONY: template-log-level
template-log-level: ${BINARY_YQ}
	@if [[ "${STAGE}" == "development" ]]; then \
      echo "Setting LOG_LEVEL env in deployment to ${LOG_LEVEL}!" ; \
      $(BINARY_YQ) -i e ".setup.env.logLevel=\"${LOG_LEVEL}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
    fi

.PHONY: template-image-pull-policy
template-image-pull-policy: $(BINARY_YQ)
	@if [[ "${STAGE}" == "development" ]]; then \
          echo "Setting pull policy to always!" ; \
          $(BINARY_YQ) -i e ".setup.imagePullPolicy=\"Always\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
    fi

.PHONY: template-dogu-registry
template-dogu-registry: $(BINARY_YQ)
	@if [[ "${STAGE}" == "development" ]]; then \
          echo "Template dogu registry!" ; \
          $(BINARY_YQ) -i e ".dogu_registry_secret.url=\"${DOGU_REGISTRY_URL}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".dogu_registry_secret.username=\"${DOGU_REGISTRY_USERNAME}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".dogu_registry_secret.password=\"${DOGU_REGISTRY_PASSWORD}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".dogu_registry_secret.urlschema=\"${DOGU_REGISTRY_URL_SCHEMA}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
    fi

.PHONY: template-docker-registry
template-docker-registry: $(BINARY_YQ)
	@if [[ "${STAGE}" == "development" ]]; then \
          echo "Template docker registry!" ; \
          $(BINARY_YQ) -i e ".container_registry_secrets[0].url=\"${DOCKER_REGISTRY_URL}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".container_registry_secrets[0].username=\"${DOCKER_REGISTRY_USERNAME}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".container_registry_secrets[0].password=\"${DOCKER_REGISTRY_PASSWORD}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
    fi

.PHONY: template-helm-registry
template-helm-registry: $(BINARY_YQ)
	@if [[ "${STAGE}" == "development" ]]; then \
          echo "Template helm registry!" ; \
          $(BINARY_YQ) -i e ".helm_registry_secret.host=\"${HELM_REGISTRY_HOST}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".helm_registry_secret.username=\"${HELM_REGISTRY_USERNAME}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".helm_registry_secret.password=\"${HELM_REGISTRY_PASSWORD}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".helm_registry_secret.schema=\"${HELM_REGISTRY_SCHEMA}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".helm_registry_secret.plainHttp=\"${HELM_REGISTRY_PLAIN_HTTP}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
          $(BINARY_YQ) -i e ".helm_registry_secret.insecureTls=\"${HELM_REGISTRY_INSECURE_TLS}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
    fi


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
	@if [[ "${RUNTIME_ENV}" == "local" ]]; then \
		echo "Cleaning in namespace $(NAMESPACE)"; \
		kubectl delete --all dogus --namespace=$(NAMESPACE) || true; \
		kubectl delete component k8s-cert-manager --namespace=$(NAMESPACE) || true; \
		kubectl delete component k8s-cert-manager-crd --namespace=$(NAMESPACE) || true; \
		kubectl delete component k8s-velero --namespace=$(NAMESPACE) || true; \
		for cmp in $$(kubectl get component --namespace=$(NAMESPACE) --output=jsonpath="{.items[*].metadata.name}"); do \
		if [[ $$cmp != *"k8s-longhorn"* ]] && [[ $$cmp != *"k8s-component-operator"* ]] && [[ $$cmp != *"k8s-component-operator-crd"* ]]; then \
				kubectl delete component $${cmp} --namespace=$(NAMESPACE); \
			fi; \
		done; \
		kubectl delete component k8s-longhorn --namespace=$(NAMESPACE) || true; \
		kubectl patch component k8s-component-operator -p '{"metadata":{"finalizers":null}}' --type=merge --namespace=$(NAMESPACE) || true; \
		kubectl patch component k8s-component-operator-crd -p '{"metadata":{"finalizers":null}}' --type=merge --namespace=$(NAMESPACE) || true; \
		helm uninstall k8s-component-operator --namespace=$(NAMESPACE) || true; \
		helm uninstall k8s-component-operator-crd --namespace=$(NAMESPACE) || true; \
		kubectl patch cm tcp-services -p '{"metadata":{"finalizers":null}}' --type=merge --namespace=$(NAMESPACE) || true; \
		kubectl patch cm udp-services -p '{"metadata":{"finalizers":null}}' --type=merge --namespace=$(NAMESPACE) || true; \
		kubectl delete statefulsets,deploy,secrets,cm,svc,sa,rolebindings,roles,clusterrolebindings,clusterroles,cronjob,pvc,pv,networkpolicy --ignore-not-found -l app=ces --namespace=$(NAMESPACE); \
		kubectl delete secrets --ignore-not-found -l name=k8s-ces-setup --namespace=$(NAMESPACE); \
		kubectl delete secret --ignore-not-found ecosystem-certificate --namespace=$(NAMESPACE); \
	fi

.PHONY: build-setup
build-setup: ${SRC} compile ## Builds the setup Go binary.

.PHONY: run
run: ## Run a setup from your host.
	go run ./main.go
