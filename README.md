# Bio-Cleaner24 REST API

## Key generation
All keys must be stored in `./keys` directory and paths can be adjusted in `auth.go` file.

### for JWT
```
$ openssl genrsa -out app.rsa 4096
$ openssl rsa -in app.rsa -pubout > app.rsa.pub
```
### for HTTPS
```
$ openssl genrsa -out server.key 4096
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3560
```