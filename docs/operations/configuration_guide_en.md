# Configuration guide

This document describes the deployment of a valid `k8s-ces-setup` configuration and explains all possible
configuration options.

## Create sample configuration

First, the configuration must be downloaded from the repository at `k8s/k8s-ces-setup-config.yaml`. The
file contains a ConfigMap with important configuration for the `k8s-ces-setup`:

```yaml
#
# The default configuration map for the ces-setup. Should always be deployed before the setup itself.
#
apiVersion: v1
kind: ConfigMap
metadata:
  name: k8s-ces-setup-config
  namespace: default
  labels:
    app: cloudogu-ecosystem
    app.kubernetes.io/name: k8s-ces-setup
data:
  k8s-ces-setup.yaml: |
    namespace: "ecosystem-0"
    logLevel: "debug"
    doguOperatorVersion: "0.0.0"
    etcdServerVersion: "0.0.0"
```

Under the `data` section the content of a `k8s-ces-setup.yaml` is defined.

## Explanation of the configuration values

### namespace

* YAML key: `namespace`
* type: `string`
* Required configuration
* Description: The namespace defines the target namespace for the Cloudogu EcoSystem to be created. This can be changed to
  be changed to any value. The namespace and all necessary components are created during the setup process.
  created.

### log_level

* YAML key: `log_level`
* Type: one of the following values `ERROR, WARN, INFO, DEBUG`
* Necessary configuration
* Description: Sets the log level of the `k8s-ces-setup` and thus how accurate the log output of the application should be.
  should be.

### dogu_operator_version

* YAML key: `dogu_operator_version`
* Type: `String` as link to the desired [Dogu Operator](http://github.com/cloudogu/k8s-dogu-operator) version
* Necessary configuration
* Description: The Dogu Operator is a central component in the EcoSystem and must be installed. The given
  link points to the version of the Dogu Operator to be installed. The link must point to a valid K8s YAML resource of the
  `k8s-dogu-operator`. This will be appended to the release of the `k8s-dogu-operator` on each release.
* Example: `TODO: Add first link when the first release is done`

### etcd_server_version

* YAML key: `etcd_server_version`
* Type: `String` as link to the desired [Etcd](http://github.com/cloudogu/k8s-etcd) version
* Necessary configuration
* Description: The Etcd is a central component in the EcoSystem and must be installed. The specified link
  points to the version of the EcoSystem Etcd to be installed. The link must point to a valid K8s YAML resource of the
  `k8s-etcd`. This is located directly in the repository under the path `manifests/etcd.yaml`.
* Example: `https://github.com/cloudogu/k8s-etcd/blob/develop/manifests/etcd.yaml`

## Deploy configuration

The created configuration can now be run via Kubectl with the following command:

```bash
kubectl apply -f k8s-ces-setup-config.yaml
```

Now the setup can be deployed. For more information about deploying the setup
[here](installation_guide_en.md).