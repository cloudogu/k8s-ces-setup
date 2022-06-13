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

### Automatic setup via setup.json

If the setup is to be performed automatically without any user interaction, this can be done using a `setup.json`.
This file contains all the configuration values required to perform the setup. How the `setup.json` can be created and
inserted into the cluster is described in ["Deployment of a setup configuration"](custom_setup_configuration_en.md).

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
wget https://raw.githubusercontent.com/cloudogu/k8s-ces-setup/develop/k8s/k8s-ces-setup.yaml
yq "(select(.kind == \"ClusterRoleBinding\").subjects[]|select(.name == \"k8s-ces-setup\")).namespace=\"your-target-namespace\"" k8s-ces-setup.yaml > k8s-ces-setup.patched.yaml

kubectl --namespace your-target-namespace apply -f k8s-ces-setup.patched.yaml
```

The k8s-ces-setup should now have started successfully in the cluster. The setup should now be reachable via the IP of the machine at
the port `30080`.

### Execute Setup

```bash
curl -I --request POST --url http://your-cluster-ip-or-fqdn:30080/api/v1/setup
```

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
