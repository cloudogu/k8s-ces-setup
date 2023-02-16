## Verwendung von selbst signierten Zertifikaten

Falls die konfigurierte Dogu-Registry ein selbst signiertes Zertifikat verwendet muss dieses über ein Secret wie bei dem
k8s-dogu-operator konfiguriert werden.

```bash
kubectl --namespace <cesNamespace> create secret generic dogu-registry-cert --from-file=dogu-registry-cert.pem=<cert_name>.pem
```

Bei einem Neustart des Controllers wird das Zertifikat nach `/etc/ssl/certs/<cert_name>.pem` gemountet und ist
für verwendete Http-Funktionen des Setups verfügbar.
