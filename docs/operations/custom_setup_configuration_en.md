# Deployment of a setup configuration

This document describes the setup configuration, its individual components and their deployment in the form of a `setup.json` file.

## Setup configuration (`setup.json`)

The setup configuration describes uniformly the data which are needed when creating a new EcoSystem.
It is possible to save a setup configuration, or parts of it, in an additional file in JSON format.
This file can be added to the 'k8s-ces-setup' to run the setup partially or completely automatically.

## Structure of a setup configuration

The setup configuration is divided by content into several sections, also called regions:

* **Naming**: Contains general configurations for the system.
* **UserBackend**: Contains configurations for the user connection.
* **AdminUser**: Contains configurations for the initial admin user in the EcoSystem.
* **Dogus**: Contains configurations for the dogus to be installed.
* **RegistryConfig**: Contains configurations which are written to the internal etcd during setup.
* **RegistryConfigEncrypted**: Contains configurations which are written encrypted to the internal Etcd during setup.

For a completely automatic setup all necessary regions must be defined in the `setup.json`.
A complete description of the individual regions and their configuration values follows in a later chapter.

## Differences to conventional `ces-setup`

The `k8s-ces-setup` differs from the [ces-setup](https://github.com/cloudogu/ces-setup) in that an EcoSystem does not run on a single VM, but within a Kubernetes cluster on multiple VMs.
A `setup.json` from the `ces-setup` can be used as setup configuration for the `k8s-ces-setup` without any problems.

**Note**: However, it should be noted that some regions/configuration values in the `k8s-ces-setup` are invalid or not yet supported. Also, the `official/nginx` dogu has been replaced by `k8s/nginx-static` and `k8s/nginx-ingress`.

### Region tokens

Since the `k8s-ces-setup` cannot configure VM's anymore, this section is omitted completely.
Properties such as `locale`, `timezone` and `keyboardLayout` must happen when the Kubernetes cluster is initialized.


### `Naming` region

The `naming` region contains configurations that affect the entire system. This includes FQDN, domain, SSL certificates and more.

Object name: _naming_
Properties:

#### useInternalIp
* Optional
* Data type: boolean
* Contents: This switch specifies whether a specific IP address should be used for an internal DNS resolution of the host. If this switch is set to `true`, then this forces a valid value in the `internalIp` field. If this field is not set, then it will be interpreted as `false` and ignored.
* Example: `"useInternalIp": true`

#### internalIp
* Optional
* Data type: String
* Contents: If and only if `userInternalIp` is true, the IP address stored here will be used for an internal DNS resolution of the host. Otherwise this field is ignored. This is especially interesting for installations with a split DNS configuration, i.e. if the instance is reachable from outside with a different IP address than from inside.
* Example: `"internalIp": "10.0.2.15"`.

The internal IP is used in `ces-setup` to write an additional entry in `etc/hosts`.
In the Kubernetes environment, this is not possible in this way and is not currently implemented.

### UserBackend section

Properties have no differences to `ces-setup`

### AdminUser region

Properties do not have any differences to `ces-setup`

### Region Dogus

Properties have no differences to the `ces-setup`

### Region RegistryConfig

Properties have no differences to the `ces-setup`

### Region RegistryConfigEncrypted

Properties have no differences to `ces-setup`
Note, however, that the key/value pairs are not set immediately in the Dogu configuration,
because the dogu operator generates the public and private key for a dogu only at the time of the dogu installation.
Therefore, the entries from the `registryConfigEncrypted` region are stored in Secrets between.
These are consumed by the dogu operator when a dogu is installed.

## Deployment of a setup configuration

If a setup configuration is available in the form of a `setup.json`, it can be spawned with the following command for the setup:

```bash
kubectl --namespace your-target-namespace create configmap k8s-ces-setup-json --from-file=setup.json
```

Now the setup can be deployed. For more information about deploying the setup
[here](installation_guide_en.md) describe.
