check-ssl-cert
==========

Check the ssl certification's expiry.

## Synopsis

```shell
check-ssl-cert --host mackerel.io --warning 30 --critical 7
```

## Options

```
Application Options:
  -H, --host=            Host name
  -p, --port=            Port number (default: 443)
  -w, --warning=days     The warning threshold in days before expiry (default: 30)
  -c, --critical=days    The critical threshold in days before expiry (default: 14)
```
