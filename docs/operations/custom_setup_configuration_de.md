# Ausbringung einer Setup-Konfiguration

Dieses Dokument beschreibt die Setup-Konfiguration, ihre einzelnen Bestandteile und ihre Ausbringung in Form einer `setup.json`-Datei.

## Setup-Konfiguration (`setup.json`)

Die Setup-Konfiguration beschreibt einheitlich die Daten, die beim Anlegen eines neuen EcoSystems benötigt werden.
Es besteht die Möglichkeit, eine Setup-Konfiguration, oder Teile davon, in einer zusätzlichen Datei im JSON-Format zu speichern.
Diese Datei kann dem `k8s-ces-setup` hinzugefügt werden, um das Setup teilweise oder ganz automatisch auszuführen.

## Aufbau einer Setup-Konfiguration

Die Setup-Konfiguration wird inhaltlich in mehrere Bereiche, auch Regionen genannt, aufgeteilt:

* **Token** - //TODO
* **Naming**: Enthält allgemeine Konfigurationen für das System.
* **UserBackend**: Enthält Konfigurationen für die Benutzer-Anbindung.
* **AdminUser**: Enthält Konfigurationen für den initialen Admin-Benutzer im EcoSystem.
* **Dogus**: Enthält Konfigurationen für die zu installierenden Dogus.
* **RegistryConfig**: Enthält Konfigurationen, welche beim Setup in den internen Etcd geschrieben werden. 
* **RegistryConfigEncrypted**: Enthält Konfigurationen, welche beim Setup verschlüsselt in den internen Etcd geschrieben werden.

Für ein komplett automatisches Setup müssen alle notwendigen Regionen in der `setup.json` definiert sein. 
Eine komplette Beschreibung der einzelnen Regionen und ihrer Konfigurationswerte folgt in einem späteren Kapitel. 

## Unterschiede zum herkömmlichen `ces-setup`

Das `k8s-ces-setup` unterscheidet sich zum [ces-setup](https://github.com/cloudogu/ces-setup) dabei, dass ein EcoSystem nicht auf einer einzelnen VM läuft, sondern innerhalb eines Kubernetes-Cluster auf mehreren VMs.
Eine `setup.json` von dem `ces-setup` kann ohne Probleme als Setup Konfiguration für das `k8s-ces-setup` benutzt werden.

**Hinweis**: Es ist jedoch zu beachten, dass einige Regionen/Konfigurationswerte im `k8s-ces-setup` hinfällig oder noch nicht unterstützt werden. 
Folgenden Regionen/Konfigurationswerte werden derzeitig von dem `k8s-ces-setup` ignoriert:

* **Token** - //TODO.
* **Region**: Enthält allgemeine Konfigurationswerte für regionale Einstellungen an der VM.
* **Projects**: Enthält Konfigurationswerte für die initiale Ausbringung von Projekten beim Setup.
* **UnixUser**: Enthält die Konfiguration für den User der VM.
* **UnattendedUpgrades**: Enthält Konfigurationen für unbeaufsichtigte Updates der VM.
* **ExtendedConfiguration**: Enthält Konfigurationswerte für spezielle Fälle.
* **SequentialDoguStart**: Konfigurationswert, um die Installation der Dogus sequenziell durchzuführen.

## Ausbringung einer Setup-Konfiguration

Wenn eine Setup-Konfiguration in Form einer `setup.json` vorliegt, kann diese mit dem folgenden Befehl für das Setup ausgebracht werden:

```bash
kubectl --namespace your-target-namespace create configmap k8s-ces-setup-json --from-file=setup.json
```

Nun kann das Setup ausgebracht werden. Für mehr Informationen zur Ausbringung des Setup sind
[hier](installation_guide_de.md) beschreiben.

## Ausführliche Beschreibung aller Regionen der Setup-Konfiguration

### Region Token

TODO

### Region Naming

Die Region `Naming` enthält Konfigurationen, die das gesamte System betreffen. Darunter fallen FQDN, Domäne, SSL-Zertifikate und mehr.

TODO

### Region UserBackend

TODO

### Region AdminUser

TODO

### Region Dogus

TODO

### Region RegistryConfig

TODO

### Region RegistryConfigEncrypted

TODO
