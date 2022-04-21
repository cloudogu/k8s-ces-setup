# Über die innere Architektur des Setups

Dieses Dokument beschreibt, den inneren Aufbau und die Funktionsweise von k8s-ces-setup. Dieses Dokument unterliegt voraussichtlich noch weiteren Bearbeitungen, da aktuell nur die Basisaktionen einer CES-Installation umgesetzt wurden. 

Hinweise zur Installation des Setups selbst liegen in der [Installationsanleitung](../operations/installation_guide_de.md) vor.

## Installationsablauf eines automatischen Setups

![Grober Installationsablauf im Setup](../images/setup-installation-workflow-overview.png)

Das automatische Setup einer CES-Instanz ganz ohne weitere Benutzerinteraktion ("unattended setup") soll so ähnlich wie möglich zum herkömmlichen CES-Setup stattfinden.

**Voraussetzungen für ein unattended Setup:**
1. [Dogu-](https://github.com/cloudogu/k8s-dogu-operator/blob/develop/docs/operations/configuring_the_dogu_registry_de.md) und [Image-Instanz credentials](https://github.com/cloudogu/k8s-dogu-operator/blob/develop/docs/operations/configuring_the_docker_registry_de.md) werden als secrets bereitgestellt
1. [Setup-Konfiguration](../operations/configuration_guide_de.md) liegt in einer `ConfigMap` im gleichen Namespace, in dem das Setup läuft
1. (noch nicht umgesetzt) ein Setup-Deskriptor `setup.json` liegt vor

**Durchführung:**
Es erfolgt in mehreren Schritten, die die Abbildung oben gut veranschaulicht:

1. Setup-Konfiguration einlesen
2. Cluster-Konfiguration einlesen
   - im Produktionsbetrieb wird diese durch das Setup-Deployment bereitgestellt
   - im Entwicklungsbetrieb kann diese auch von lokalen Kube-Configs ausgelesen werden
3. Dogu- und Image-Credentials auslesen
4. neuen Namespace (gemäß Setup-Konfiguration) anlegen
5. ausgelesene Credentials in den neuen Namespace kopieren 
6. etcd-Server in den neuen Namespace installieren
7. etcd-Client in den neuen Namespace installieren
8. Dogu-Operator in den neuen Namespace installieren
9. (noch nicht umgesetzt) Dogus (gemäß `setup.json` installieren)

## (Unstructured) YAML-Ressourcen auf die K8s-API anwenden

Die Funktionalität hinter dem zentralen Struct `core.k8sApplyClient` basiert auf der Beschreibung von Ymmt2005[s Blogartikel](https://ymmt2005.hatenablog.com/entry/2020/04/14/An_example_of_using_dynamic_client_of_k8s.io/client-go#Mapping-between-GVK-and-GVR).

`core.k8sApplyClient` funktioniert ähnlich wie `kubectl` in dem es:
- K8s-Ressourcen in YAML-Form ausliest
- für die K8s-API verständlich transformiert
- und an die K8s-API

Die K8s-API verwendet REST-Calls auf Basis von JSON. YAML ist daher grundsätzlich nicht kompatibel und muss transformiert werden. Wegen der komplexen Natur von K8s-Ressourcen (viele unterschiedliche Ressourcen in teilweise vielen unterschiedlichen Versionen) ist eine manuelle Transformation in JSON oder mit sonst üblichen typisierten Mechanismen (`clientSet.CoreV1().Namespaces().Create(...)`) nicht machbar. Ebenso wenig lässt sich der Quellcode von `kubectl` sinnvoll wiederverwenden, da dessen innere Struktur so stark auf ein Kommandozeilentool ausgelegt wurde, dass der zentrale Kern der YAML-Transformation nicht sinnvoll übernommen werden kann.

Stattdessen wird ein dynamischer Mechanismus verwendet, der K8s-Ressourcen in `unstructured.Unstructured` einliest. Daraus wird die entsprechende REST-API ermitteln, die durch [Server Side Apply](https://kubernetes.io/docs/reference/using-api/api-concepts/#server-side-apply) dann auf den Cluster angewendet wird. Durch die Natur des `PATCH`-Verbs im REST-Call kann so eine noch nicht existierende Ressource angelegt oder eine bereits existierende Ressource aktualisiert werden.

**Die Schritte als Übersicht:**

1. Vorbereitung: REST-Mapper herstellen, um das GVR zu finden
2. Vorbereitung: Dynamic client herstellen
3. YAML nach `unstructured.Unstructured`-Objekt parsen
4. GVR durch GVK und REST-Mapper ermitteln
5. REST-Interface für das GVR ermitteln
6. Objekt nach JSON umwandeln
7. Ressource durch PATCH anlegen oder aktualisieren (deutet der API ein `Server Side Apply` an)
