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
    component_operator_chart: "k8s/k8s-component-operator:0.0.2"
    components:
      "k8s/k8s-etcd": "3.5.7-4"
      "k8s/k8s-dogu-operator": "0.35.0"
      "k8s/k8s-service-discovery": "0.13.0"
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

### component_operator_chart

* YAML-Key: `component_operator_chart`
* Typ: `String` als HelmChart-Bezeichner des [Komponenten-Operator](http://github.com/cloudogu/k8s-component-operator) (inkl. Namespace und Version)
* Notwendige Konfiguration
* Beschreibung: Der Komponenten-Operator ist eine zentrale Komponente im EcoSystem und muss installiert werden. Das angegebene HelmChart gibt die zu installierende Version des Komponenten-Operators an.
* Beispiel: `k8s/k8s-component-operator:0.0.2`

> **Hinweis:** als Version kann "latest" angegeben werden um die höchste, verfügbare Version des Komponenten-Operators zu verwenden.

### components

* YAML-Key: `components`
* Typ: `Map` von zu installierenden CES-Komponenten und der jeweiligen Version
* Notwendige Konfiguration
* Beschreibung: Das Setup installiert alle angegebenen CES Komponenten mit Hilfe des [Komponenten-Operators](http://github.com/cloudogu/k8s-component-operator). Folgende Komponenten werden u.a. benötigt: [Dogu Operator](http://github.com/cloudogu/k8s-dogu-operator), [Service Discovery](http://github.com/cloudogu/k8s-service-discovery), [Etcd](http://github.com/cloudogu/k8s-etcd)
* Beispiel: 
  ```yaml
    components:
      "k8s/k8s-etcd": "3.5.7-4"
      "k8s/k8s-dogu-operator": "0.35.0"
      "k8s/k8s-service-discovery": "0.13.0"
  ```

> **Hinweis:** als Version kann "latest" angegeben werden um die höchste, verfügbare Version der jeweiligen Komponente zu verwenden.

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
    * `dogu`: Diese Phase findet nach der Erzeugung von K8s Dogu-Ressourcen statt
    * `component`: Diese Phase findet nach der Erzeugung von K8s-Cloudogu-EcoSystem-Komponenten-Ressourcen statt
  * **zu ändernde Ressourcen**: Um Kubernetes-Ressourcen im Cluster-Namespace adressieren zu können, muss in Kubernetes-Syntax die jeweilige Ressource beschrieben werden Siehe hierzu auch [Objects In Kubernetes (engl.)](https://kubernetes.io/docs/concepts/overview/working-with-objects/). Ferner wird bei Ressourcen mit Namespace-Bezug der [Namespace](#beispiel-konfiguration-anlegen) verwendet, in dem das Setup des EcoSystems konfiguriert wurde.
    * `apiVersion`: Die Gruppe (optional bei K8s-Core-Ressourcen) und Version der Kubernetes-Ressource. 
    * `kind`: Die Art der Kubernetes-Ressource
    * `name`: Der konkrete Name der einzelnen Ressource
  * **JSON-Patch**: Eine Liste von einem oder mehreren JSON-Patches, die auf die Ressource angewendet werden sollen, siehe hierzu [JSON-Patch RFC 6902](https://datatracker.ietf.org/doc/html/rfc6902). Es werden diese Operationen unterstützt:
    * `add` zum Hinzufügen neuer Werte
       * für diese Operation muss ein `value`-Feld mit dem neuen Wert existieren
    * `replace` zum Ersetzen bestehender Werte mit neuen Werten
       * für diese Operation muss ein `value`-Feld mit dem neuen Wert existieren
    * `remove` zum Löschen bestehender Werte
       * diese Operation akzeptiert kein `value`-Feld 

Beispiel: 

```yaml
resource_patches:
  - phase: dogu
    resource:
      # the usual notation of Kubernetes resources is used here.
      apiVersion: k8s.cloudogu.com/v1
      kind: dogu
      name: nexus
    patches:
      # A YAML representation of JSON is used here, which is easier to write. Direct JSON is also allowed
      - op: add
        path: /spec/additionalIngressAnnotations
        value:
          nginx.ingress.kubernetes.io/proxy-body-size: "0"
      - op: replace
        path: /spec/resources
        value:
          dataVolumeSize: 5Gi
      - op: delete
        path: /spec/fieldWithATypo
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