# SSL

## Ablage

Das SSL-Zertifikat befindet sich in der Registry unter den folgenden Pfaden:
- `config/_global/certificate/key`
- `config/_global/certificate/server.crt`
- `config/_global/certificate/server.key`

## Ein SSL-Zertifikat erzeugen

Wenn die `setup.json` einen `self-signed` Zertifikatstyp angibt, ist die Erzeugung Teil des Setup-Prozesses.
Andernfalls schreibt das Setup einfach das angegebene externe Zertifikat in die Registry.
Ein `self-signed` Zertifikat ist für 365 Tage gültig. Es kann jedoch erneuert werden, indem man ein neues Zertifikat mit der folgenden Anfrage erzeugt:

```bash
curl -I --request POST --url http://fqdn:30080/api/v1/ssl?days=<days> 
```
