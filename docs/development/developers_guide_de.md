# Entwickleranleitung

Dieses Dokument Informationen um die Entwicklung an dem `k8s-ces-setup` zu unterstützen.

## Lokale Entwicklung

2. Die lokale Entwicklung am Setup kann mit `go run .` gestartet werden.

## Makefile-Targets

Der Befehl `make help` gibt alle verfügbaren Targets und deren Beschreibungen in der Kommandozeile aus.

Damit auch die Makefiles bezüglich des Clusters funktionieren muss der Root Path der Entwicklungsumgebung in den 
Makefiles unter der Umgebungsvariable `K8S_CLUSTER_ROOT` eingetragen werden.