# Entwickleranleitung

Dieses Dokument Informationen um die Entwicklung an dem `k8s-ces-setup` zu unterstützen.

## Lokale Entwicklung

Die lokale Entwicklung am Setup kann mit `go run .` gestartet werden.

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
