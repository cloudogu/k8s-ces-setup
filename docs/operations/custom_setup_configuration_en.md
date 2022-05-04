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

**Note**: However, it should be noted that some regions/configuration values in the `k8s-ces-setup` are invalid or not yet supported.
The following regions/configuration values are currently ignored by the `k8s-ces-setup`:

* **Token** - //TODO.
* **Region**: Contains general configuration values for regional settings on the VM.
* **Projects**: Contains configuration values for initial deployment of projects during setup.
* **UnixUser**: Contains configuration for the user of the VM.
* **UnattendedUpgrades**: Contains configurations for unattended updates to the VM.
* **ExtendedConfiguration**: Contains configuration values for special cases.
* **SequentialDoguStart**: Configuration value to perform the Dogus installation sequentially.

## Deployment of a setup configuration

If a setup configuration is available in the form of a `setup.json`, it can be spawned with the following command for the setup:

```bash
kubectl --namespace your-target-namespace create configmap k8s-ces-setup-json --from-file=setup.json
```

Now the setup can be deployed. For more information about deploying the setup
[here](installation_guide_en.md) describe.

## Detailed description of all regions of the setup configuration

### Region Token

TODO

### Region Naming

The 'Naming' region contains configurations that affect the entire system. This includes FQDN, domain, SSL certificates and more.

TODO

### Region UserBackend

TODO

### Region AdminUser

TODO

### Region Dogus

TODO

### Region RegistryConfig

TODO

### Region RegistryConfigEncrypted

TODO
