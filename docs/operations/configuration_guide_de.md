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
    log_level: "DEBUG"
    dogu_operator_url: https://github.com/cloudogu/k8s-dogu-operator/releases/download/v0.8.0/k8s-dogu-operator_0.8.0.yaml
    service_discovery_url: https://github.com/cloudogu/k8s-service-discovery/releases/download/v0.3.0/k8s-service-discovery_0.3.0.yaml
    static_webserver_url: https://github.com/cloudogu/k8s-static-webserver/releases/download/v0.1.0/k8s-static_webserver_0.1.0.yaml
    etcd_server_url: https://raw.githubusercontent.com/cloudogu/k8s-etcd/develop/manifests/etcd.yaml
    etcd_client_image_repo: bitnami/etcd:3.5.2-debian-10-r0
    key_provider: pkcs1v15
    remote_registry_url_schema: default
```

Unter dem Abschnitt `data`-Abschnitt wird der Inhalt einer `k8s-ces-setup.yaml` definiert.

## Erklärung der Konfigurationswerte

### log_level

* YAML-Key: `log_level`
* Typ: einer der folgenden Werte `ERROR, WARN, INFO, DEBUG`
* Notwendig Konfiguration
* Beschreibung: Setzt das Log Level des `k8s-ces-setup` und somit wie genau die Log-Ausgaben der Applikation sein sollen.

### dogu_operator_version

* YAML-Key: `dogu_operator_version`
* Typ: `String` als Link zu der gewünschten [Dogu Operator](http://github.com/cloudogu/k8s-dogu-operator) Version
* Notwendig Konfiguration
* Beschreibung: Der Dogu Operator ist eine zentrale Komponente im EcoSystem und muss installiert werden. Der angegebene Link zeigt auf die zu installierende Version des Dogu Operators. Der Link muss auf eine valide K8s-YAML-Ressource des `k8s-dogu-operator` zeigen. Diese wird bei jeder Veröffentlichung an das Release des `k8s-dogu-operator` gehängt.
* Beispiel: `https://github.com/cloudogu/k8s-dogu-operator/releases/download/v0.2.0/k8s-dogu-operator_0.2.0.yaml`

### service_discovery_url

* YAML-Key: `service_discovery_url`
* Typ: `String` als Link zu der gewünschten [Service Discovery](http://github.com/cloudogu/k8s-service-discovery) Version
* Notwendig Konfiguration
* Beschreibung: Die Service Discovery ist eine zentrale Komponente im EcoSystem und muss installiert werden. Der angegebene Link zeigt auf die zu installierende Version der Service Discovery. Der Link muss auf eine valide K8s-YAML-Ressource der `k8s-service-discovery` zeigen. Diese wird bei jeder Veröffentlichung an das Release der `k8s-service-discovery` gehängt.
* Beispiel: `https://github.com/cloudogu/k8s-service-discovery/releases/download/v0.1.0/k8s-service-discovery_0.1.0.yaml`

### static_webserver_url

* YAML-Key: `static_webserver_url`
* Typ: `String` als Link zu der gewünschten Version des [Statischen Webserver](http://github.com/cloudogu/k8s-static-webserver)
* Notwendig Konfiguration
* Beschreibung: Der Statischen Webserver ist eine zentrale Komponente im EcoSystem und muss installiert werden. Der angegebene Link zeigt auf die zu installierende Version des Statischen Webserver. Der Link muss auf eine valide K8s-YAML-Ressource der `static_webserver_url` zeigen. Diese wird bei jeder Veröffentlichung an das Release der `static_webserver_url` gehängt.
* Beispiel: `https://github.com/cloudogu/k8s-static-webserver/releases/download/v0.1.0/k8s-static-webserver_0.1.0.yaml`

### etcd_server_url

* YAML-Key: `etcd_server_url`
* Typ: `String` als Link zu der gewünschten [Etcd](http://github.com/cloudogu/k8s-etcd) Version
* Notwendig Konfiguration
* Beschreibung: Der Etcd ist eine zentrale Komponente im EcoSystem und muss installiert werden. Der angegebene Link zeigt auf die zu installierende Version des EcoSystem-Etcd. Der Link muss auf eine valide K8s-YAML-Ressource des `k8s-etcd` zeigen. Diese liegt direkt im Repository unter dem Pfad `manifests/etcd.yaml`.
* Beispiel: `https://github.com/cloudogu/k8s-etcd/blob/develop/manifests/etcd.yaml`

### etcd_client_image_repo

* YAML-Key: `etcd_client_image_repo`
* Typ: `String` als Name zum gewünschten [Etcd-Client](https://artifacthub.io/packages/helm/bitnami/etcd) Image.
* Notwendig Konfiguration
* Beschreibung: Der Etcd-Client ist eine Komponente im EcoSystem welche die Kommunikation mit dem Etcd-Server vereinfacht. Der Eintrag muss auf ein valides Image von `bitnami/etcd` sein.
* Beispiel: `bitnami/etcd:3.5.2-debian-10-r0`

### key_provider

* YAML-Key: `key_provider`
* Typ: einer der folgenden Werte `pkcs1v15, oaesp`
* Notwendig Konfiguration
* Beschreibung: Setzt den verwendeten Key-Provider des Ecosystems und beeinflusst so die zu verschlüsselnde Registry-Werte.
* Beispiel: `pkcs1v15`

### remote_registry_url_schema

* YAML-Key: `remote_registry_url_schema`
* Typ: einer der folgenden Werte `default, index`
* Notwendig Konfiguration
* Beschreibung: Setzt das URLSchema der Remote-Registry.
* Beispiel: `default` in normalen Umgebungen, `index` in gespiegelten Umgebungen

## Konfiguration ausbringen

Die erstellte Konfiguration kann nun via Kubectl mit dem folgenden Befehl ausgeführt werden:

```bash
kubectl apply -f k8s-ces-setup-config.yaml
```

Nun kann das Setup ausgebracht werden. Für mehr Informationen zur Ausbringung des Setup sind
[hier](installation_guide_de.md) beschreiben.