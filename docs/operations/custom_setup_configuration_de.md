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

### Region Token

Da das `k8s-ces-setup` keine VM's mehr konfigurieren kann entfällt dieser Abschnitt komplett.
Eigenschaften wie `locale`, `timezone` und `keyboardLayout` müssen bei der Initialisierung des Kubernetes-Cluster geschehen.


### Bereich "Naming"

Die Region `Naming` enthält Konfigurationen, die das gesamte System betreffen. Darunter fallen FQDN, Domäne, SSL-Zertifikate und mehr.

Objektname: _naming_
Eigenschaften:

#### useInternalIp
* Optional
* Datentyp: boolean
* Inhalt: Dieser Schalter gibt an, ob eine spezifische IP-Adresse für eine interne DNS-Auflösung des Hosts verwendet werden soll. Wenn dieser Schalter auf `true` gesetzt wird, dann erzwingt dies, einen gültigen Wert im Feld `internalIp`. Wenn dieses Feld nicht gesetzt wurde, dann wird es mit `false` interpretiert und ignoriert.
* Beispiel: `"useInternalIp": true`

#### internalIp
* Optional
* Datentyp: String
* Inhalt: Wenn und nur wenn `userInternalIp` wahr ist, wird die hier hinterlegte IP-Adresse für eine interne DNS-Auflösung des Hosts verwendet. Ansonsten wird dieses Feld ignoriert. Dies ist besonders für Installationen mit einer Split-DNS-Konfiguration interessant, d. h. wenn die Instanz von außen mit einer anderen IP-Adresse erreichbar ist, als von innen.
* Beispiel: `"internalIp": "10.0.2.15"`

Die interne IP wird im `ces-setup` dazu verwendet um einen zusätzlichen Eintrag in `etc/hosts` zu schreiben.
Im Kubernetes-Umfeld ist dies so nicht möglich und zurzeit nicht implementiert.

### Bereich "UserBackend"

Eigenschaften besitzen keine Unterschiede zum `ces-setup`

### Region AdminUser

Eigenschaften besitzen keine Unterschiede zum `ces-setup`

### Region Dogus

Eigenschaften besitzen keine Unterschiede zum `ces-setup`

### Region RegistryConfig

Eigenschaften besitzen keine Unterschiede zum `ces-setup`

### Region RegistryConfigEncrypted

Eigenschaften besitzen keine Unterschiede zum `ces-setup`.
Zu beachten ist allerdings, dass die Schlüssel/Wert-Paare nicht sofort in der Dogu-Konfiguration gesetzt werden, 
weil der Dogu-Operator erst zum Zeitpunkt der Dogu-Installation den Public- und Private-Key für ein Dogu erzeugt.
Daher werden die Einträge aus der Region `registryConfigEncrypted` in Secrets zwischen gespeichert.
Diese werden bei der Installation eines Dogus von dem Dogu Operator konsumiert.

## Ausbringung einer Setup-Konfiguration

Wenn eine Setup-Konfiguration in Form einer `setup.json` vorliegt, kann diese mit dem folgenden Befehl für das Setup ausgebracht werden:

```bash
kubectl --namespace your-target-namespace create configmap k8s-ces-setup-json --from-file=setup.json
```

Nun kann das Setup ausgebracht werden. Für mehr Informationen zur Ausbringung des Setup sind
[hier](installation_guide_de.md) beschreiben.
