# Installation guide

This document describes all necessary steps to install the 'k8s-ces-setup'.

## Prerequisites

1. a running K8s cluster exists.
2. `kubectl` has been installed and has been configured for the existing K8s cluster.

## Installation from GitHub

### Deploy configuration

The `k8s-ces-setup` needs a configuration for the installation. This must be provided in the form of a ConfigMap before the
before installing the `k8s-ces-setup`. More information about deployment and the individual
configuration options is described [here](configuration_guide_en.md).

### Deploy setup

The installation from GitHub requires the installation YAML, which contains all the required K8s resources. This is located
in the repository under `k8s/k8s-ces-setup.yaml`. The installation looks like this with `kubectl`:

```
kubectl apply -f https://github.com/cloudogu/k8s-ces-setup/blob/develop/k8s/k8s-ces-setup.yaml
```

The k8s-ces-setup should now have started successfully in the cluster. The setup should now be reachable via the IP of the machine at
the port `30080`.