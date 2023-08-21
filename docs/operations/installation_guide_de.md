# Installationsanleitung

Dieses Dokument beschreibt alle notwendigen Schritte um das `k8s-ces-setup` zu installieren.

## Voraussetzungen

1. Ein laufendes K8s-Cluster ist vorhanden.
2. `kubectl` wurde installiert und wurde für das vorhandene K8s-Cluster konfiguriert.
3. `helm` wurde installiert.

## Installation mit Helm

### Automatisches Setup via setup.json

Soll das Setup automatisch ohne Anwenderinteraktion durchgeführt werden, kann dies mithilfe einer `setup.json` geschehen.
Diese enthält alle nötigen Konfigurationswerte zur Durchführung des Setups. Wie die `setup.json` erstellt und in den
Cluster eingebracht werden kann, ist in ["Ausbringung einer Setup-Konfiguration"](custom_setup_configuration_de.md) beschrieben.

### Setup ausbringen

Die Installation mit Helm erfordert die Konfiguration der `values.yaml`. Ein minimales Beispiel sieht folgendermaßen aus:

```yaml
docker_registry_secret:
  url: https://registry.cloudogu.com
  username: "your-ces-instance-id"
  password: "your-ces-instance-password"

dogu_registry_secret:
  url: https://dogu.cloudogu.com/api/v2/dogus
  username: "your-ces-instance-id"
  password: "your-ces-instance-password"

helm_registry_secret:
  url: https://registry.cloudogu.com
  username: "your-ces-instance-id"
  password: "your-ces-instance-password"

# Example test setup.json
#setup_json:
#  {
#    "naming": {
#      "fqdn": "",
#      "domain": "k3ces.local",
#      "certificateType": "selfsigned",
#      "relayHost": "yourrelayhost.com",
#      "useInternalIp": false,
#      "internalIp": ""
#      "completed": true,
#    },
#    "dogus": {
#      "defaultDogu": "redmine",
#      "install": [
#        "official/ldap",
#        "official/postfix",
#        "k8s/nginx-static",
#        "k8s/nginx-ingress",
#        "official/cas",
#        "official/postgresql",
#        "official/redmine",
#      ],
#      "completed": true
#    },
#    "admin": {
#      "username": "admin",
#      "mail": "admin@admin.admin",
#      "password": "adminpw",
#      "adminGroup": "cesAdmin",
#      "adminMember": true,
#      "sendWelcomeMail": false,
#      "completed": true
#    },
#    "userBackend": {
#      "dsType": "embedded",
#      "server": "",
#      "attributeID": "uid",
#      "attributeGivenName": "",
#      "attributeSurname": "",
#      "attributeFullname": "cn",
#      "attributeMail": "mail",
#      "attributeGroup": "memberOf",
#      "baseDN": "",
#      "searchFilter": "(objectClass=person)",
#      "connectionDN": "",
#      "password": "",
#      "host": "ldap",
#      "port": "389",
#      "loginID": "",
#      "loginPassword": "",
#      "encryption": "",
#      "groupBaseDN": "",
#      "groupSearchFilter": "",
#      "groupAttributeName": "",
#      "groupAttributeDescription": "",
#      "groupAttributeMember": "",
#      "completed": true
#    }
#  }
```

> Für weitere Konfigurationen wie z.B. Versionen der Operatoren siehe [values.yaml](https://github.com/cloudogu/k8s-ces-setup/blob/feature/59_helm_release/k8s/helm/values.yaml).

### Setup installieren

- `helm registry login registry.cloudogu.com --username "your-ces-instance-id" --password "your-ces-instance-password"`
- `helm upgrade -i -f values.yaml k8s-ces-setup oci//:registry.cloudogu.com/k8s/k8s-ces-setup `

### Setup ausführen

- `kubectl port-forward service/k8s-ces-setup 30080:8080`
- `curl -I --request POST --url http://localhost:30080/api/v1/setup`

### Status des Setups

Für die Präsentation des Zustands existiert eine ConfigMap `k8s-setup-config` mit dem Data-Key
`state`. Mögliche werte sind `installing, installed`. Falls der Wert `installing` vor dem Setup-Prozess gesetzt sind, bricht ein
Start des Setups sofort ab.

`kubectl --namespace your-target-namespace describe configmap k8s-setup-config`

Falls der Wert `installed` gesetzt ist, ist das Setup bereit aus dem Cluster gelöscht zu werden.

### Cleanup des Setups

Mit dem Setup wird ein CronJob `k8s-ces-setup-finisher` ausgeliefert der periodisch (default: 1 Minute) prüft, ob das Setup erfolgreich durchlaufen ist.
Tritt dieser Fall ein werden alle Ressourcen mit dem Label `app.kubernetes.io/name=k8s-ces-setup` gelöscht.
Zusätzlich werden Konfigurationen wie z.B. die `setup.json` und der CronJob selbst entfernt. Cluster bezogene Objekte werden nicht gelöscht.

Da der CronJob nicht seine eigene Rolle löschen kann, muss diese als einzige Ressource manuell entfernt werden:
`kubectl delete role k8s-ces-setup-finisher`
