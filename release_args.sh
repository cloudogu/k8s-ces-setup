#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail
# This script is automatically called by the automatic git flow release process. It is responsible to change the
# version of the image in the K8s deployment resource `k8s/k8s-ces-setup.yaml` to the newest one.
valuesYaml=k8s/helm/values.yaml
patchTplYaml=k8s/helm/component-patch-tpl.yaml

update_versions_modify_files() {
  newReleaseVersion="${1}"

  yq -i ".setup.image.tag=\"${newReleaseVersion}\"" "${valuesYaml}"

  setupImage="cloudogu/k8s-ces-setup:${newReleaseVersion}"
  yq -i ".values.images.k8sCesSetup=\"${setupImage}\"" "${patchTplYaml}"

  kubectlImage=$(yq ".kubectl_image" "${valuesYaml}")
  yq -i ".values.images.kubectl=\"${kubectlImage}\"" "${patchTplYaml}"
}

update_versions_stage_modified_files() {
  git add "${valuesYaml}" "${patchTplYaml}"
}