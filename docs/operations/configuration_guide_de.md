# Konfigurationsanleitung

Dieses Dokument beschreibt die Ausbringung einer validen `k8s-ces-setup`-Konfiguration und erklärt alle möglichen
Konfigurationsoptionen.

## Beispiel-Konfiguration anlegen

Als Erstes muss die Konfiguration aus dem Repository unter `k8s/k8s-ces-setup-config.yaml` heruntergeladen werden. Die
Datei enthält eine ConfigMap mit wichtiger Konfiguration für das `k8s-ces-setup`:

```yaml
#
# The default configuration map for the ces-setup. Should always be deployed before the setup itself.
#
apiVersion: v1
kind: ConfigMap
metadata:
  name: k8s-ces-setup-config
  namespace: default
  labels:
    app: cloudogu-ecosystem
    app.kubernetes.io/name: k8s-ces-setup
data:
  k8s-ces-setup.yaml: |
    namespace: "ecosystem-0"
    logLevel: "debug"
    doguOperatorVersion: "0.0.0"
    etcdServerVersion: "0.0.0"
```

Unter dem Abschnitt `data`-Abschnitt wird der Inhalt einer `k8s-ces-setup.yaml` definiert.

## Erklärung der Konfigurationswerte

### namespace

* YAML-Key: `namespace`
* Typ: `String`
* Notwendig Konfiguration
* Beschreibung: Der Namespace definiert den Ziel-Namespace für das zu erstellende Cloudogu EcoSystem. Dieser kann zu
  einem beliebigen Wert geändert werden. Der Namespace und alle notwendigen Komponenten werden im Verlauf des Setups
  angelegt.

### log_level

* YAML-Key: `log_level`
* Typ: einer der folgenden Werte `ERROR, WARN, INFO, DEBUG`
* Notwendig Konfiguration
* Beschreibung: Setzt das Log Level des `k8s-ces-setup` und somit wie genau die Log-Ausgaben der Applikation sein
  sollen.

### dogu_operator_version

* YAML-Key: `dogu_operator_version`
* Typ: `String` als Link zu der gewünschten [Dogu Operator](http://github.com/cloudogu/k8s-dogu-operator) Version
* Notwendig Konfiguration
* Beschreibung: Der Dogu Operator ist eine zentrale Komponente im EcoSystem und muss installiert werden. Der angegebene
  Link zeigt auf die zu installierende Version des Dogu Operators. Der Link muss auf eine valide K8s-YAML-Ressource des
  `k8s-dogu-operator` zeigen. Diese wird bei jeder Veröffentlichung an das Release des `k8s-dogu-operator` gehängt.
* Beispiel: `TODO: Add first link when the first release is done`

### etcd_server_version

* YAML-Key: `etcd_server_version`
* Typ: `String` als Link zu der gewünschten [Etcd](http://github.com/cloudogu/k8s-etcd) Version
* Notwendig Konfiguration
* Beschreibung: Der Etcd ist eine zentrale Komponente im EcoSystem und muss installiert werden. Der angegebene Link
  zeigt auf die zu installierende Version des EcoSystem-Etcd. Der Link muss auf eine valide K8s-YAML-Ressource des
  `k8s-etcd` zeigen. Diese liegt direkt im Repository unter dem Pfad `manifests/etcd.yaml`.
* Beispiel: `https://github.com/cloudogu/k8s-etcd/blob/develop/manifests/etcd.yaml`

## Konfiguration ausbringen

Die erstellte Konfiguration kann nun via Kubectl mit dem folgenden Befehl ausgeführt werden:

```bash
kubectl apply -f k8s-ces-setup-config.yaml
```

Nun kann das Setup ausgebracht werden. Für mehr Informationen zur Ausbringung des Setup sind
[hier](installation_guide_de.md) beschreiben.