# Entwickleranleitung

Dieses Dokument Informationen um die Entwicklung an dem `k8s-ces-setup` zu unterstützen.

## Lokale Entwicklung

Zuerst sollten Entwicklungsdateien angelegt werden, die anstelle der Cluster-Werte verwendet werden sollen:

* `k8s/dev-config/k8s-ces-setup.yaml`: [setup-config](../operations/configuration_guide_de.md)
* `k8s/dev-config/setup.json`: [custom-setup-config](../operations/custom_setup_configuration_de.md)

### Installation des Ces-Setups im lokalen Cluster

Damit das ces-setup im lokalen Cluster ausgeführt und getestet werden kann, müssen einige Dinge beachtet werden.
Zuerst sollten alle vorhandenen Dogus, Komponenten, etc. aus dem System entfernt werden. Dazu kann
der Befehl `make k8s-clean` verwendet werden.
Damit anschließend das Ces-Setup installiert werden kann, muss vorher noch eine kleine Änderung an der 
k8s/helm/values.yaml durchgeführt werden.
Der folgende Teil muss einkommentiert werden, andernfalls kann das Setup nicht durchgeführt werden:
```
  # k8s-longhorn:
  #   version: latest
  #   helmRepositoryNamespace: k8s
  #   deployNamespace: longhorn-system
```
Anschließend kann mit `make helm-apply` das Ces-Setup installiert werden. Es wird dann automatisch durchgeführt.


### Ausführung mit `go run` oder einer IDE

- die lokale Entwicklung am Setup kann mit `STAGE=development go run .` gestartet werden
- Ausführung und Debugging in IDEs wie IntelliJ ist möglich
  - allerdings sollte hierbei die Umgebungsvariable `STAGE` ebenfalls nicht vergessen werden

## Makefile-Targets

Der Befehl `make help` gibt alle verfügbaren Targets und deren Beschreibungen in der Kommandozeile aus.

## Debugging

Es ist möglich, mit einem deployten Setup zu interagieren:

```bash
# Setup-Zustand prüfen
curl --request GET --url http://192.168.56.2:30080/api/v1/health
{"status":"healthy","version":"0.0.0"}

# Namespace laut Setup configuration map anlegen
curl -I --request POST --url http://192.168.56.2:30080/api/v1/setup
```

## Pre-Setup-Zustand herstellen

Manchmal ist es notwendig, die Zeit wieder auf Anfang zurückzudrehen, z. B. um Installationsroutinen zu überprüfen. 
Dies lässt sich mit dem Make-Target `k8s-clean` erreichen (auf den **aktuellen Namespace** achten):

```bash
# delete all dogus & components and all the resources directly created by the setup
make k8s-clean

# eventuell noch fälschlich ausgebrachte Ressourcen manuell löschen
...
```

## Cleanup des Setups

Wenn in der Entwicklung neue Kubernetes-Ressourcen erstellt werden, müssen ggf. diese von dem Cleanup-Task berücksichtigt werden.
Dazu ist das Skript der Configmap `k8s-ces-setup-cleanup-script` in der `k8s-ces-setup.yaml` zu bearbeiten.