#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail
# This script is automatically called by the automatic git flow release process. It is responsible to change the
# version of the image in the K8s deployment resource `k8s/k8s-ces-setup.yaml` to the newest one.

update_versions_modify_files() {
  newReleaseVersion="${1}"
  newImage="cloudogu/k8s-ces-setup:${newReleaseVersion}"
  k8sCesSetupYaml=k8s/k8s-ces-setup.yaml

  yq e "(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.name == \"k8s-ces-setup\")).image=\"${newImage}\"" \
    ${k8sCesSetupYaml} > tmpfile

  mv tmpfile "${k8sCesSetupYaml}"
}

update_versions_stage_modified_files() {
  k8sCesSetupYaml=k8s/k8s-ces-setup.yaml

  git add "${k8sCesSetupYaml}"
}