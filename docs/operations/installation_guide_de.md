# Installationsanleitung

Dieses Dokument beschreibt alle notwendigen Schritte um das `k8s-ces-setup` zu installieren.

## Voraussetzungen

1. Ein laufendes K8s-Cluster ist vorhanden.
2. `kubectl` wurde installiert und wurde für das vorhandene K8s-Cluster konfiguriert. 

## Installation von GitHub

### Konfiguration ausbringen

Das `k8s-ces-setup` benötigt eine Konfiguration für die Installation. Diese muss in Form einer ConfigMap vor der 
Installation des `k8s-ces-setup` ausgebracht werden. Mehr Information zur Ausbringung und zu den einzelnen 
Konfigurationsoptionen wird [hier](configuration_guide_de.md) beschrieben.

### Setup ausbringen

Die Installation von GitHub erfordert die Installations-YAML, die alle benötigten K8s-Ressourcen enthält. Diese liegt
im Repository unter `k8s/k8s-ces-setup.yaml`. Die Installation sieht mit `kubectl` folgendermaßen aus:

```
kubectl apply -f https://github.com/cloudogu/k8s-ces-setup/blob/develop/k8s/k8s-ces-setup.yaml
```

Das k8s-ces-setup sollte nun erfolgreich im Cluster gestartet sein. Das Setup sollte nun über die IP der Maschine unter
dem Port `30080` erreichbar sein.