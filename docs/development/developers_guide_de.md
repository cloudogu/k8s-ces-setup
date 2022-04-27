# Entwickleranleitung

Dieses Dokument Informationen um die Entwicklung an dem `k8s-ces-setup` zu unterstützen.

## Lokale Entwicklung

Zuerst sollten Entwicklungsdateien angelegt werden, die anstelle der Cluster-Werte verwendet werden sollen:

Dogu-Operator-Resource:
- eine passende YAML-Datei (z. B. `dev-dogu-operator.yaml`) unter `k8s/dev-resources/` ablegen
- `make serve-local-yaml` liefert alle Ressourcen in dem Verzeichnis aus
  - Test: [http://localhost:9876/](http://localhost:9876/)
  - ein DNS-/Host-Alias ist hilfreich, um vom lokalen K8s-Cluster mit diesem HTTP-Server zu kommunizieren 
  - das Target benötigt Python3

`k8s/dev-config/k8s-ces-setup.yaml`:
- `namespace` legt fest, in welchem Namespace das Cloudogu EcoSystem installiert werden soll
- `dogu_operator_url` legt die Dogu-Operator-Resource fest
  - z. B. `http://192.168.56.1:9876/dev-dogu-operator.yaml` (siehe oben)

### Ausführung mit `go run` oder einer IDE

- die lokale Entwicklung am Setup kann mit `STAGE=development go run .` gestartet werden
- Ausführung und Debugging in IDEs wie IntelliJ ist möglich
  - allerdings sollte hierbei die Umgebungsvariable `STAGE` ebenfalls nicht vergessen werden

## Makefile-Targets

Der Befehl `make help` gibt alle verfügbaren Targets und deren Beschreibungen in der Kommandozeile aus.

Damit auch die Makefiles bezüglich des Clusters funktionieren muss der Root Path der Entwicklungsumgebung in den 
Makefiles unter der Umgebungsvariable `K8S_CLUSTER_ROOT` eingetragen werden.

## Debugging

Es ist möglich, mit einem deployten Setup zu interagieren:

```bash
# Setup-Zustand prüfen
curl --request GET --url http://192.168.56.2:30080/api/v1/health
{"status":"healthy","version":"0.0.0"}

# Namespace laut Setup configuration map anlegen
curl -I --request POST --url http://192.168.56.2:30080/api/v1/setup
```

## Pre-Setup-Zustang herstellen

Manchmal ist es notwendig, die Zeit wieder auf Anfang zurückzudrehen, z. B. um Installationsroutinen zu überprüfen. Dies lässt sich mit den folgenden Befehlen erreichen (auf den **aktuellen Namespace** achten):

```bash
# delete the resources directly created by the setup
make k8s-delete
# löscht Zielnamespace und alle darin Namespaced Ressourcen (pods, deployments, secrets, usw.)
kubectl delete ns your-namespace
# löscht CRD, sodass diese initial mit dem Dogu-Operator eingespielt werden kann
kubectl delete crd dogus.k8s.cloudogu.com
# löscht clusterroles/bindings aus setup-Installationen
kubectl delete clusterroles k8s-dogu-operator-metrics-reader ingress-nginx
kubectl delete clusterrolebindings ingress-nginx
# eventuell noch fälschlich ausgebrachte Ressourcen manuell löschen
...
```
