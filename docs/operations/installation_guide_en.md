# Installation guide

This document describes all necessary steps to install the 'k8s-ces-setup'.

## Prerequisites

1. a running K8s cluster exists.
2. `kubectl` has been installed and has been configured for the existing K8s cluster.
3. `helm` has been installed.

## Installation with Helm

### Automatic setup via setup.json

If the setup is to be performed automatically without any user interaction, this can be done using a `setup.json`.
This file contains all the configuration values required to perform the setup. How the `setup.json` can be created and
inserted into the cluster is described in ["Deployment of a setup configuration"](custom_setup_configuration_en.md).

### Deploy setup

The installation with helm requires the configuration of the `values.yaml`. Passwords for the registries need to be 
base64-encoded ([see here](configuration_guide_en.md#base64-encoding-tips)).
A minimal example would be:

```yaml
docker_registry_secret:
  url: https://registry.cloudogu.com
  username: "your-ces-instance-id"
  password: "eW91ci1jZXMtaW5zdGFuY2UtcGFzc3dvcmQ=" # Base64 encoded password

dogu_registry_secret:
  url: https://dogu.cloudogu.com/api/v2/dogus
  username: "your-ces-instance-id"
  password: "eW91ci1jZXMtaW5zdGFuY2UtcGFzc3dvcmQ=" # Base64 encoded password

helm_registry_secret:
  host: https://registry.cloudogu.com
  schema: oci
  plainHttp: "false"
  username: "your-ces-instance-id"
  password: "eW91ci1jZXMtaW5zdGFuY2UtcGFzc3dvcmQ=" # Base64 encoded password

component_operator_crd_chart: "k8s/k8s-component-operator-crd:latest"
component_operator_chart: "k8s/k8s-component-operator:latest"

components:
  k8s-etcd:
    version: latest
    helmRepositoryNamespace: k8s
  k8s-dogu-operator:
    version: latest
    helmRepositoryNamespace: k8s
  k8s-dogu-operator-crd:
    version: latest
    helmRepositoryNamespace: k8s
  k8s-service-discovery:
    version: latest
    helmRepositoryNamespace: k8s

# Example test setup.json
#setup_json:
#  {
#    "naming": {
#      "fqdn": "",
#      "domain": "k3ces.local",
#      "certificateType": "selfsigned",
#      "relayHost": "yourrelayhost.com",
#      "useInternalIp": false,
#      "internalIp": ""
#      "completed": true,
#    },
#    "dogus": {
#      "defaultDogu": "redmine",
#      "install": [
#        "official/ldap",
#        "official/postfix",
#        "k8s/nginx-static",
#        "k8s/nginx-ingress",
#        "official/cas",
#        "official/postgresql",
#        "official/redmine",
#      ],
#      "completed": true
#    },
#    "admin": {
#      "username": "admin",
#      "mail": "admin@admin.admin",
#      "password": "adminpw",
#      "adminGroup": "cesAdmin",
#      "adminMember": true,
#      "sendWelcomeMail": false,
#      "completed": true
#    },
#    "userBackend": {
#      "dsType": "embedded",
#      "server": "",
#      "attributeID": "uid",
#      "attributeGivenName": "",
#      "attributeSurname": "",
#      "attributeFullname": "cn",
#      "attributeMail": "mail",
#      "attributeGroup": "memberOf",
#      "baseDN": "",
#      "searchFilter": "(objectClass=person)",
#      "connectionDN": "",
#      "password": "",
#      "host": "ldap",
#      "port": "389",
#      "loginID": "",
#      "loginPassword": "",
#      "encryption": "",
#      "groupBaseDN": "",
#      "groupSearchFilter": "",
#      "groupAttributeName": "",
#      "groupAttributeDescription": "",
#      "groupAttributeMember": "",
#      "completed": true
#    }
#  }
```
<!-- markdown-link-check-disable-next-line -->
> For more configuration options like operator versions see [values.yaml](https://github.com/cloudogu/k8s-ces-setup/blob/develop/k8s/helm/values.yaml).

### Install the Setup

- `helm registry login registry.cloudogu.com --username "your-ces-instance-id" --password "your-ces-instance-password"`
- `helm upgrade -i -f values.yaml k8s-ces-setup oci//:registry.cloudogu.com/k8s/k8s-ces-setup `

### Execute Setup

- `kubectl port-forward service/k8s-ces-setup 30080:8080`
- `curl -I --request POST --url http://localhost:30080/api/v1/setup`

### Status of the setup

For the presentation of the state there is a ConfigMap `k8s-setup-config` with the data key
`state`. Possible values are `installing, installed`. If these values are set before the setup process, a start of the setup
start of the setup will abort immediately.

`kubectl --namespace your-target-namespace describe configmap k8s-setup-config`

### Cleanup of the setup

A cron job `k8s-ces-setup-finisher` is delivered with the setup which periodically (default: 1 minute) checks whether the setup has run successfully.
If this occurs, all resources with the label `app.kubernetes.io/name=k8s-ces-setup` are deleted.
Additionally, configurations such as `setup.json` and the cron job itself are removed. Cluster scoped objects are not deleted.

Since the cron job cannot delete its own role, it is the only resource that must be removed manually:
`kubectl delete role k8s-ces-setup-finisher`.
