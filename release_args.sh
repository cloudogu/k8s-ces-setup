#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

# this function will be sourced from release.sh and be called from release_functions.sh
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