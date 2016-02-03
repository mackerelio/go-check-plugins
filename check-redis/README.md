# check-redis

## Description

Checks for Redis

## Sub Commands

- reachable
- slave

## check-redis reachable

Checks if Redis is reachable.

### Setting
* need a host and port pair, or a socket
```
[plugin.checks.redis_reachable]
command = "/path/to/check-redis reachable [--host=127.0.0.1] [--port=6379] [--timeout=5] [--socket=<unix socket>]
```

## check-redis slave

Checks Redis slave status.

### Setting
* need a host and port pair, or a socket
```
[plugin.checks.redis_slave]
command = "/path/to/check-redis slave [--host=127.0.0.1] [--port=6379] [--timeout=5] [--socket=<unix socket>]
```
