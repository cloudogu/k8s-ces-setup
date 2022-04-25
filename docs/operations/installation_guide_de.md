# Installationsanleitung

Dieses Dokument beschreibt alle notwendigen Schritte um das `k8s-ces-setup` zu installieren.

## Voraussetzungen

1. Ein laufendes K8s-Cluster ist vorhanden.
2. `kubectl` wurde installiert und wurde für das vorhandene K8s-Cluster konfiguriert. 

## Installation von GitHub

### Konfiguration ausbringen

Das `k8s-ces-setup` benötigt eine Konfiguration für die Installation. Diese muss in Form einer ConfigMap vor der 
Installation des `k8s-ces-setup` ausgebracht werden. Mehr Information zur Ausbringung und zu den einzelnen 
Konfigurationsoptionen wird [im Configuration-Guide](configuration_guide_de.md) beschrieben.

### Setup ausbringen

Die Installation von GitHub erfordert die Installations-YAML, die alle benötigten K8s-Ressourcen enthält. Diese liegt
im Repository unter `k8s/k8s-ces-setup.yaml`. Die Installation sieht mit `kubectl` folgendermaßen aus:

```bash
kubectl create ns your-target-namespace
kubectl create secret docker-registry k8s-dogu-operator-dogu-registry --namespace=ecosystem --docker-server=registry.cloudogu.com --docker-username="your-ces-instance-id" --docker-password="your-ces-instance-password"
kubectl create secret generic k8s-dogu-operator-docker-registry --namespace=ecosystem --from-literal=username="your-ces-instance-id" --from-literal=password="your-ces-instance-password"

kubectl apply -f https://github.com/cloudogu/k8s-ces-setup/blob/develop/k8s/k8s-ces-setup.yaml
```

Das k8s-ces-setup sollte nun erfolgreich im Cluster gestartet sein. Das Setup sollte nun über die IP der Maschine unter
dem Port `30080` erreichbar sein.

### Setup ausführen

```bash
curl -I --request POST --url http://your-cluster-ip-or-fqdn:30080/api/v1/setup
```