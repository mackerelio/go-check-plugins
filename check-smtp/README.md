# check-smtp

## Description

Check for SMTP connection.

## Setting

```
[plugin.checks.smtp]
command = "/path/to/check-smtp -H smtp.example.com -p 25 -w 3 -c 5 -t 10"
```

## Options

```
  -H, --host=         Hostname (default: localhost)
  -p, --port=         Port (default: 25)
  -F, --fqdn=         FQDN used for HELO
  -s, --smtps         Use SMTP over TLS
  -S, --starttls      Use STARTTLS
  -A, --authmech=     SMTP AUTH Authentication Mechanisms(only PLAIN supported)
  -U, --authuser=     SMTP AUTH username
  -P, --authpassword= SMTP AUTH password
  -w, --warning=      Warning threshold (sec) (default: 3)
  -c, --critical=     Critical threshold (sec) (default: 5)
  -t, --timeout=      Timeout (sec) (default: 10)
```
