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


