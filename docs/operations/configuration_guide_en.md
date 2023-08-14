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
  namespace: ecosystem
  labels:
    app: cloudogu-ecosystem
    app.kubernetes.io/name: k8s-ces-setup
data:
  k8s-ces-setup.yaml: |
    log_level: "DEBUG"
    dogu_operator_url: https://dogu.cloudogu.com/api/v1/k8s/k8s/k8s-dogu-operator
    service_discovery_url: https://dogu.cloudogu.com/api/v1/k8s/k8s/k8s-service-discovery
    etcd_server_url: https://raw.githubusercontent.com/cloudogu/k8s-etcd/develop/manifests/etcd.yaml
    etcd_client_image_repo: bitnami/etcd:3.5.2-debian-10-r0
    key_provider: pkcs1v15
    resource_patches:
    - phase: dogu
      resource:
        apiVersion: k8s.cloudogu.com/v1
        kind: Dogu
        name: nexus
      patches:
        - op: add
          path: /spec/resources
          value:
            dataVolumeSize: 5Gi
```

Under the `data` section the content of a `k8s-ces-setup.yaml` is defined.
The `namespace` entry must correspond to the namespace in the cluster where the CES is to be installed.

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

### etcd_server_url

* YAML key: `etcd_server_url`
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

### resource_patches

* YAML key: `resource_patches`
* Type: list of patch objects
* Optional configuration
* Description: list of patch objects that are applied to Kubernetes resources at different stages of setup, e.g., to apply user- or environment-specific changes. These patch objects consist of three components: Setup Phase, Resource to Change, and JSON Patch.
   * **Setup Phases**: These phases currently exist:
      * `loadbalancer`: this phase occurs after the Kubernetes load balancer service is created.
      * `dogu`: This phase takes place after the creation of K8s dogu resources.
      * `component`: This phase takes place after the creation of K8s Cloudogu EcoSystem component resources.
   * **resources to modify**: To be able to address Kubernetes resources in the cluster namespace, the respective resource must be described in Kubernetes syntax. See also [Objects In Kubernetes](https://kubernetes.io/docs/concepts/overview/working-with-objects/). Furthermore, resources with namespace reference use the [namespace](#create-sample-configuration) in which the EcoSystem setup was configured.
      * `apiVersion`: The group (optional for K8s core resources) and version of the Kubernetes resource.
      * `kind`: The type of Kubernetes resource.
      * `name`: The specific name of the individual resource.
   * **JSON patch**: A list of one or more JSON patches to apply to the resource, see [JSON patch RFC 6902](https://datatracker.ietf.org/doc/html/rfc6902). These operations are supported:
      * `add` to add new values
         * for this operation, a `value` field with the new value is required
      * `replace` to replace existing values with new values.
         * for this operation, a `value` field with the new value is required
      * `remove` to delete existing values
        * for this operation, any `value` definition must be absent 

Example:

```yaml
resource_patches:
  - phase: dogu
    resource:
# the usual notation of Kubernetes resources is used here.
      apiVersion: k8s.cloudogu.com/v1
      kind: dogu
      name: nexus
    patches:
# A YAML representation of JSON is used here, which is easier to write. Direct JSON is also allowed
      - op: add
        path: /spec/additionalIngressAnnotations
        value:
          nginx.ingress.kubernetes.io/proxy-body-size: "0"
      - op: replace
        path: /spec/resources
        value:
          dataVolumeSize: 5Gi
      - op: delete
        path: /spec/fieldWithATypo
```

#### Notes on JSON Patches

`value` fields in JSON patches must form key-value pairs.

When a JSON patch needs to add an empty object as a key value (like below in the `myKey` example), this notation is used:
```yaml
resource_patches:
# ...
    patches:
      - op: add
        path: /path/to/resourcefield
        value:
          myKey: {}
```

If a JSON patch path references fields that do not exist, the Kubernetes API cannot create them recursively. Instead, the missing fields must be configured in separate patches.

```yaml
resource_patches:
# ...
    patches:
# creates "key", which probably does not exist yet
      - op: add
        path: /spec/key
        value: {}
# now the key "anotherKey" can be added to "key"
      - op: add
        path: /spec/key/anotherKey
        value:
          response: 42
```

## Deploy configuration

The created configuration can now be run via Kubectl with the following command:

```bash
kubectl apply -f k8s-ces-setup-config.yaml
```

Now the setup can be deployed. For more information about deploying the setup
[here](installation_guide_en.md).

## Configuration of the index-URL scheme.

If you want the k8s-ces-setup to install Dogus from a Dogu registry with index-URL scheme, you have to specify this in the
cluster secret `k8s-dogu-operator-dogu-registry`. This secret is created during the k8s-dogu-operator configuration,
see https://github.com/cloudogu/k8s-dogu-operator/blob/develop/docs/operations/configuring_the_dogu_registry_en.md.
The secret has to contain the key `urlschema`, which should be set to `index`. If this key is not present
or not set to `index`, the `default` URL scheme is used.
