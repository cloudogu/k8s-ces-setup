# Installation guide

This document describes all necessary steps to install the 'k8s-ces-setup'.

## Prerequisites

1. a running K8s cluster exists.
2. `kubectl` has been installed and has been configured for the existing K8s cluster.

## Installation from GitHub

### Deploy configuration

The `k8s-ces-setup` needs a configuration for the installation. This must be provided in the form of a ConfigMap before the
before installing the `k8s-ces-setup`. More information about deployment and the individual
configuration options describes [the configuration guide](configuration_guide_en.md).

### Deploy setup

The installation from GitHub requires the installation YAML, which contains all the required K8s resources. This is located
in the repository under `k8s/k8s-ces-setup.yaml`. The installation looks like this with `kubectl`:

```bash
kubectl create ns your-target-namespace
kubectl create secret generic k8s-dogu-operator-dogu-registry \
    --namespace=your-target-namespace \
    --from-literal=endpoint="https://dogu.cloudogu.com/api/v2/dogus" \
    --from-literal=username="your-ces-instance-id" \
    --from-literal=password="your-ces-instance-password"
kubectl create secret docker-registry k8s-dogu-operator-docker-registry \
    --namespace=your-target-namespace \
    --docker-server=registry.cloudogu.com \
    --docker-username="your-ces-instance-id" \
    --docker-password="your-ces-instance-password"

# note: the setup resource must be modified with your-target-namespace
wget https://github.com/cloudogu/k8s-ces-setup/blob/develop/k8s/k8s-ces-setup.yaml
yq "(select(.kind == \"ClusterRoleBinding\").subjects[]|select(.name == \"k8s-ces-setup\")).namespace=\"your-target-namespace\"" k8s-ces-setup.yaml > k8s-ces-setup.patched.yaml

kubectl --namespace your-target-namespace apply -f k8s-ces-setup.patched.yaml
```

The k8s-ces-setup should now have started successfully in the cluster. The setup should now be reachable via the IP of the machine at
the port `30080`.

### Execute Setup

```bash
curl -I --request POST --url http://your-cluster-ip-or-fqdn:30080/api/v1/setup
```