check-http
==========

Check http.

## Synopsis

```shell
check-http -u http://example.com
```

To override status
```shell
check-http -s 404=ok -u http://example.com
check-http -s 200-404=ok -u http://example.com
```

To change request destination
```shell
check-http --connect-to=example.com:443:127.0.0.1:8080 https://example.com # will request to 127.0.0.1:8000 but AS example.com:443
check-http --connect-to=:443:127.0.0.1:8080 https://example.com # empty host matches ANY host
check-http --connect-to=example.com::127.0.0.1:8080 https://example.com # empty port matches ANY port
```
