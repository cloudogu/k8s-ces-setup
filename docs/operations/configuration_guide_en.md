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
    log_level: "DEBUG"
    dogu_operator_url: https://github.com/cloudogu/k8s-dogu-operator/releases/download/v0.2.0/k8s-dogu-operator_0.2.0.yaml
    service_discovery_url: https://github.com/cloudogu/k8s-service-discovery/releases/download/v0.1.0/k8s-service-discovery_0.1.0.yaml
    etcd_server_url: https://raw.githubusercontent.com/cloudogu/k8s-etcd/develop/manifests/etcd.yaml
    etcd_client_image_repo: bitnami/etcd:3.5.2-debian-10-r0
```

Under the `data` section the content of a `k8s-ces-setup.yaml` is defined.

## Explanation of the configuration values

### log_level

* YAML key: `log_level`
* Type: one of the following values `ERROR, WARN, INFO, DEBUG`
* Necessary configuration
* Description: Sets the log level of the `k8s-ces-setup` and thus how accurate the log output of the application should be.

### dogu_operator_version

* YAML key: `dogu_operator_version`
* Type: `String` as link to the desired [Dogu Operator](http://github.com/cloudogu/k8s-dogu-operator) version
* Necessary configuration
* Description: The Dogu Operator is a central component in the EcoSystem and must be installed. The given link points to the version of the Dogu Operator to be installed. The link must point to a valid K8s YAML resource of the `k8s-dogu-operator`. This will be appended to the release of the `k8s-dogu-operator` on each release.
* Example: `https://github.com/cloudogu/k8s-dogu-operator/releases/download/v0.2.0/k8s-dogu-operator_0.2.0.yaml`

### service_discovery_url

* YAML key: `service_discovery_url`
* Type: `String` as link to the desired [Service Discovery](http://github.com/cloudogu/k8s-service-discovery) version.
* Necessary configuration
* Description: Service Discovery is a central component in EcoSystem and must be installed. The specified link points to the version of Service Discovery to be installed. The link must point to a valid K8s YAML resource of the `k8s-service-discovery`. This will be appended to the release of the `k8s-service-discovery` on each release.
* Example: `https://github.com/cloudogu/k8s-service-discovery/releases/download/v0.1.0/k8s-service-discovery_0.1.0.yaml`

### etcd_server_version

* YAML key: `etcd_server_version`
* Type: `String` as link to the desired [Etcd](http://github.com/cloudogu/k8s-etcd) version
* Necessary configuration
* Description: The Etcd is a central component in the EcoSystem and must be installed. The specified link points to the version of the EcoSystem Etcd to be installed. The link must point to a valid K8s YAML resource of the `k8s-etcd`. This is located directly in the repository under the path `manifests/etcd.yaml`.
* Example: `https://github.com/cloudogu/k8s-etcd/blob/develop/manifests/etcd.yaml`

### etcd_client_image_repo

* YAML key: `etcd_client_image_repo`
* Type: `String` as name to desired [Etcd-Client](https://artifacthub.io/packages/helm/bitnami/etcd) image.
* Necessary configuration
* Description: The Etcd-Client is a component in the EcoSystem which simplifies the communication with the Etcd-Server. The entry must be on a valid image of `bitnami/etcd`.
* Example: `bitnami/etcd:3.5.2-debian-10-r0`

### key_provider

* YAML key: `key_provider`
* Type: one of the following values `pkcs1v15, oaesp`
* Required configuration
* Description: Sets the used key provider of the ecosystem and thus influences the registry values to be encrypted.
* Example: `pkcs1v15`

### remote_registry_url_schema

* YAML key: `remote_registry_url_schema`
* Type: one of the following values `default, index`.
* Required Configuration
* Description: Sets the URLSchema of the remote registry.
* Example: `default` in normal environments, `index` in mirrored environments.

## Deploy configuration

The created configuration can now be run via Kubectl with the following command:

```bash
kubectl apply -f k8s-ces-setup-config.yaml
```

Now the setup can be deployed. For more information about deploying the setup
[here](installation_guide_en.md).