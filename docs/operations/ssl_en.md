# Use of a Self-Signed SSL Certificate

## Location

The SSL certificate is located in the registry under the following paths:
- `config/_global/certificate/key`
- `config/_global/certificate/server.crt`
- `config/_global/certificate/server.key`

## Generate an SSL certificate

If the `setup.json` specify an `self-signed` certificate type the generation is part of the setup process.
Otherwise, the setup just writes the given external certificate to the registry.
A self-signed certificate is valid for 365 days. However, it can be renewed by generating a new one with following request:

```bash
curl -I --request POST --url http://fqdn:30080/api/v1/ssl?days=<days> 
```
