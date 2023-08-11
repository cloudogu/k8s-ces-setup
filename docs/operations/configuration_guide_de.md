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
  namespace: ecosystem
  labels:
    app: cloudogu-ecosystem
    app.kubernetes.io/name: k8s-ces-setup
data:
  k8s-ces-setup.yaml: |
    log_level: "DEBUG"
    dogu_operator_url: https://dogu.cloudogu.com/api/v1/k8s/k8s/k8s-dogu-operator
    service_discovery_url: https://dogu.cloudogu.com/api/v1/k8s/k8s/k8s-service-discovery
    etcd_server_url: https://raw.githubusercontent.com/cloudogu/k8s-etcd/develop/manifests/etcd.yaml
    etcd_client_image_repo: bitnami/etcd:3.5.2-debian-10-r0
    key_provider: pkcs1v15
    resource_patches:
    - phase: dogu
      resource:
        apiVersion: k8s.cloudogu.com/v1
        kind: Dogu
        name: nexus
      patches:
        - op: add
          path: /spec/resources
          value:
            dataVolumeSize: 5Gi
```

Unter dem Abschnitt `data`-Abschnitt wird der Inhalt einer `k8s-ces-setup.yaml` definiert.
Der Eintrag `namespace` muss dem Namespace im Cluster entsprechen, in den das CES installiert werden soll.

## Erklärung der Konfigurationswerte

### log_level

* YAML-Key: `log_level`
* Typ: einer der folgenden Werte `ERROR, WARN, INFO, DEBUG`
* Notwendige Konfiguration
* Beschreibung: Setzt das Log Level des `k8s-ces-setup` und somit wie genau die Log-Ausgaben der Applikation sein sollen.

### dogu_operator_version

* YAML-Key: `dogu_operator_version`
* Typ: `String` als Link zu der gewünschten [Dogu Operator](http://github.com/cloudogu/k8s-dogu-operator) Version
* Notwendige Konfiguration
* Beschreibung: Der Dogu Operator ist eine zentrale Komponente im EcoSystem und muss installiert werden. Der angegebene Link zeigt auf die zu installierende Version des Dogu Operators. Der Link muss auf eine valide K8s-YAML-Ressource des `k8s-dogu-operator` zeigen. Diese wird bei jeder Veröffentlichung an das Release des `k8s-dogu-operator` gehängt.
* Beispiel: `https://github.com/cloudogu/k8s-dogu-operator/releases/download/v0.2.0/k8s-dogu-operator_0.2.0.yaml`

### service_discovery_url

* YAML-Key: `service_discovery_url`
* Typ: `String` als Link zu der gewünschten [Service Discovery](http://github.com/cloudogu/k8s-service-discovery) Version
* Notwendige Konfiguration
* Beschreibung: Die Service Discovery ist eine zentrale Komponente im EcoSystem und muss installiert werden. Der angegebene Link zeigt auf die zu installierende Version der Service Discovery. Der Link muss auf eine valide K8s-YAML-Ressource der `k8s-service-discovery` zeigen. Diese wird bei jeder Veröffentlichung an das Release der `k8s-service-discovery` gehängt.
* Beispiel: `https://github.com/cloudogu/k8s-service-discovery/releases/download/v0.1.0/k8s-service-discovery_0.1.0.yaml`

### etcd_server_url

* YAML-Key: `etcd_server_url`
* Typ: `String` als Link zu der gewünschten [Etcd](http://github.com/cloudogu/k8s-etcd) Version
* Notwendige Konfiguration
* Beschreibung: Der Etcd ist eine zentrale Komponente im EcoSystem und muss installiert werden. Der angegebene Link zeigt auf die zu installierende Version des EcoSystem-Etcd. Der Link muss auf eine valide K8s-YAML-Ressource des `k8s-etcd` zeigen. Diese liegt direkt im Repository unter dem Pfad `manifests/etcd.yaml`.
* Beispiel: `https://github.com/cloudogu/k8s-etcd/blob/develop/manifests/etcd.yaml`

### etcd_client_image_repo

* YAML-Key: `etcd_client_image_repo`
* Typ: `String` als Name zum gewünschten [Etcd-Client](https://artifacthub.io/packages/helm/bitnami/etcd) Image.
* Notwendige Konfiguration
* Beschreibung: Der Etcd-Client ist eine Komponente im EcoSystem welche die Kommunikation mit dem Etcd-Server vereinfacht. Der Eintrag muss auf ein valides Image von `bitnami/etcd` sein.
* Beispiel: `bitnami/etcd:3.5.2-debian-10-r0`

### key_provider

* YAML-Key: `key_provider`
* Typ: einer der folgenden Werte `pkcs1v15, oaesp`
* Notwendige Konfiguration
* Beschreibung: Setzt den verwendeten Key-Provider des Ecosystems und beeinflusst so die zu verschlüsselnde Registry-Werte.
* Beispiel: `pkcs1v15`

### resource_patches

* YAML-Key: `resource_patches`
* Typ: Liste von Patch-Objekten
* Optionale Konfiguration
* Beschreibung: Liste von Patch-Objekten, die zu unterschiedlichen Phasen des Setups auf Kubernetes-Ressourcen angewendet werden, z. B. um benutzer- oder umgebungsspezifische Änderungen auszubringen. Diese Patch-Objekte bestehen aus drei Bestandteilen: Setup-Phase, zu ändernde Resource und JSON-Patch
  * **Setup-Phasen**: Diese Phasen existieren aktuell:
    * `loadbalancer`: Diese Phase findet nach der Erzeugung des Kubernetes Load-Balancer-Services statt
      * Patches in dieser Phase werden nur ausgeführt, wenn die FQDN in der [Setup.json](custom_setup_configuration_de.md#Bereich-Naming) leer bzw. auf den IP-Adressenplatzhalter `<<ip>>` gesetzt wurde.
    * `dogu`: Diese Phase findet nach der Erzeugung von K8s Dogu-Ressourcen statt
    * `component`: Diese Phase findet nach der Erzeugung von K8s-Cloudogu-EcoSystem-Komponenten-Ressourcen statt
  * **zu ändernde Ressourcen**: Um Kubernetes-Ressourcen im Cluster-Namespace adressieren zu können, muss in Kubernetes-Syntax die jeweilige Ressource beschrieben werden Siehe hierzu auch [Objects In Kubernetes (engl.)](https://kubernetes.io/docs/concepts/overview/working-with-objects/). Ferner wird bei Ressourcen mit Namespace-Bezug der [Namespace](#Beispiel-Konfiguration-anlegen) verwendet, in dem das Setup des EcoSystems konfiguriert wurde.
    * `apiVersion`: Die Gruppe (optional bei K8s-Core-Ressourcen) und Version der Kubernetes-Ressource. 
    * `kind`: Die Art der Kubernetes-Ressource
    * `name`: Der konkrete Name der einzelnen Ressource
  * **JSON-Patch**: Eine Liste von einem oder mehreren JSON-Patches, die auf die Ressource angewendet werden sollen, siehe hierzu [JSON-Patch RFC 6902](https://datatracker.ietf.org/doc/html/rfc6902). Es werden diese Operationen unterstützt:
    * `add` zum Hinzufügen neuer Werte
    * `remove` zum Löschen bestehender Werte
    * `replace` zum Ersetzen bestehender Werte mit neuen Werten

Beispiel: 

```yaml
resource_patches:
  - phase: dogu
    resource:
# hier wird die übliche Schreibweise von Kubernetes-Ressourcen verwendet
      apiVersion: k8s.cloudogu.com/v1
      kind: Dogu
      name: nexus
    patches:
# Hier wird eine YAML-Repräsentation von JSON verwendet, die leichter schreibbar ist. Direktes JSON ist ebenso erlaubt
      - op: add
        path: /spec/additionalIngressAnnotations
        value:
          nginx.ingress.kubernetes.io/proxy-body-size: "0"
      - op: add
        path: /spec/resources
        value:
          dataVolumeSize: 5Gi
```

#### Hinweise zu JSON-Patches

`value`-Felder in JSON-Patches müssen Schlüssel-Wertpaare bilden.

Wenn ein JSON-Patch ein leeres Objekt als Wert für einen Schlüssel (im Beispiel `meinKey`) hinzugefügt werden soll, wird diese Notation verwendet:
```yaml
resource_patches:
# ...
    patches:
      - op: add
        path: /pfad/zu/ressourcefeld
        value:
          meinKey: {}
```

Wenn ein JSON-Patch-Pfad Felder referenziert, die nicht existieren, so kann sie die Kubernetes-API nicht rekursiv anlegen. Stattdessen müssen die fehlenden Felder in separaten Patches konfiguriert werden.

```yaml
resource_patches:
# ...
    patches:
# legt den Schlüssel "key" an, der wohl noch nicht existiert
      - op: add
        path: /spec/key
        value: {}
# nun kann der Schlüssel "nochEinKey" in "key" hinzugefügt werden
      - op: add
        path: /spec/key/nochEinKey
        value:
          antwort: 42
```

## Konfiguration ausbringen

Die erstellte Konfiguration kann nun via Kubectl mit dem folgenden Befehl ausgeführt werden:

```bash
kubectl apply -f k8s-ces-setup-config.yaml
```

Nun kann das Setup ausgebracht werden. Für mehr Informationen zur Ausbringung des Setup sind
[hier](installation_guide_de.md) beschreiben.

## Konfiguration des index-URL-Schemas

Soll das k8s-ces-setup Dogus aus einer Dogu-Registry mit index-URL-Schema installieren, muss dies
im Cluster-Secret `k8s-dogu-operator-dogu-registry` hinterlegt werden. Dieses Secret wird im Umfeld des k8s-dogu-operators
angelegt, siehe https://github.com/cloudogu/k8s-dogu-operator/blob/develop/docs/operations/configuring_the_dogu_registry_de.md.
Das Secret muss den Key `urlschema` enthalten, welcher auf `index` gesetzt sein muss. Ist dieser Key nicht vorhanden
oder nicht auf `index` gesetzt, wird das `default`-URL-Schema benutzt.