## Using self-signed certificates

If the configured Dogu registry uses a self-signed certificate, you must configure it using a Secret.

```bash
kubectl --namespace <cesNamespace> create secret generic dogu-registry-cert --from-file=dogu-registry-cert.pem=<cert_name>.pem
```

When the controller is restarted, the certificates are mounted to `/etc/ssl/certs/<cert_name>.pem` and are
available for used Https functions of the controller.
