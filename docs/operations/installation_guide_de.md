# Installationsanleitung

Dieses Dokument beschreibt alle notwendigen Schritte um das `k8s-ces-setup` zu installieren.

## Voraussetzungen

1. Ein laufendes K8s-Cluster ist vorhanden.
2. `kubectl` wurde installiert und wurde für das vorhandene K8s-Cluster konfiguriert.

## Installation von GitHub

### Konfiguration ausbringen

Das `k8s-ces-setup` benötigt eine Konfiguration für die Installation. Diese muss in Form einer ConfigMap vor der
Installation des `k8s-ces-setup` ausgebracht werden. Mehr Information zur Ausbringung und zu den einzelnen
Konfigurationsoptionen wird [im Configuration-Guide](configuration_guide_de.md) beschrieben.

### Automatisches Setup via setup.json

Soll das Setup automatisch ohne Anwenderinteraktion durchgeführt werden, kann dies mithilfe einer `setup.json` geschehen.
Diese enthält alle nötigen Konfigurationswerte zur Durchführung des Setups. Wie die `setup.json` erstellt und in den
Cluster eingebracht werden kann, ist in ["Ausbringung einer Setup-Konfiguration"](custom_setup_configuration_de.md) beschrieben.

### Setup ausbringen

Die Installation von GitHub erfordert die Installations-YAML, die alle benötigten K8s-Ressourcen enthält. Diese liegt im
Repository unter `k8s/k8s-ces-setup.yaml`. Die Installation sieht mit `kubectl` folgendermaßen aus:

```bash
kubectl create ns your-target-namespace
kubectl create secret generic k8s-dogu-operator-dogu-registry \
    --namespace=your-target-namespace \
    --from-literal=endpoint="https://dogu.cloudogu.com/api/v2/dogus" \
    --from-literal=username="your-ces-instance-id" \
    --from-literal=password="your-ces-instance-password"
kubectl create secret docker-registry k8s-dogu-operator-docker-registry \
    --namespace=your-target-namespace \
    --docker-server=registry.cloudogu.com \
    --docker-username="your-ces-instance-id" \
    --docker-password="your-ces-instance-password"
kubectl create configmap component-operator-helm-repository \
    --from-literal=endpoint="https://registry.cloudogu.com"
kubectl create secret generic component-operator-helm-registry \
    --from-literal=config.json="{\"auths\": {\"https://registry.cloudogu.com\": {\"auth\": \"$(printf "%s:%s" "your-ces-instance-id" "your-ces-instance-password" | base64)\"}}}"

# Hinweis: Die setup-Ressource muss mit dem passenden Namespace (hier: your-target-namespace) angepasst werden
wget https://raw.githubusercontent.com/cloudogu/k8s-ces-setup/develop/k8s/k8s-ces-setup.yaml
yq "(select(.kind == \"ClusterRoleBinding\").subjects[]|select(.name == \"k8s-ces-setup\")).namespace=\"your-target-namespace\"" k8s-ces-setup.yaml > k8s-ces-setup.patched.yaml

kubectl --namespace your-target-namespace apply -f k8s-ces-setup.patched.yaml
```

Das k8s-ces-setup sollte nun erfolgreich im Cluster gestartet sein. Das Setup sollte nun über die IP der Maschine unter
dem Port `30080` erreichbar sein.

### Setup ausführen

```bash
curl -I --request POST --url http://your-cluster-ip-or-fqdn:30080/api/v1/setup
```

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
