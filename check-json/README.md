check-json
==========

Check json.

## Synopsis

```shell
check-http -u http://example.com
```

To override status
```shell
check-json -s 404=ok -u http://example.com
check-json -s 200-404=ok -u http://example.com
```

To change request destination
```shell
check-json --connect-to=example.com:443:127.0.0.1:8080 https://example.com # will request to 127.0.0.1:8000 but AS example.com:443
check-json --connect-to=:443:127.0.0.1:8080 https://example.com # empty host1 matches ANY host
check-json --connect-to=example.com::127.0.0.1:8080 https://example.com # empty port1 matches ANY port
check-json --connect-to=localhost:443::8080 https://localhost # empty host2 means unchanged, therefore will request to localhost:8080 AS localhost:443
check-json --connect-to=example.com:443:127.0.0.1: https://example.com # empty port2 means unchanged, therefore will request to 127.0.0.1:443
```
