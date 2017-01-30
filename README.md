# okonomiyaki
tftp reverse proxy(backend http)

## How to User

```bash
$ ./okonomiyaki -h
Usage of ./okonomiyaki:
  -backend string
        backend server (default "http://localhost")
  -listen string
        listen address port (default ":69")
```

```bash
# Server
$ ./okonomiyaki -backend http://localhost -listen :1069
backend = http://localhost
listen  = :1069
2017/02/01 02:20:03 RRQ from 127.0.0.1:49712: index.html
2017/02/01 02:20:03 RRQ Complete from 127.0.0.1:49712: index.html
```

```bash
# Client
$ curl tftp://localhost:1069/index.html
Hello world!
```
