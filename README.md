# Go Preview

A preview generation system for NAS.

# TODO

- [ ] Support SMB protocol
- [ ] Support NFS protocol

# Build

## Local

A MySQL server should be running

```shell
export GOVIEW_DATABASE_URL="root:password@tcp(127.0.0.1:3306)/goview?charset=utf8mb4&parseTime=True&loc=Local"
go run .
```

### More environment variables

are in [env/env.go](env/env.go)

## Docker

### Build

```shell
docker build --build-arg http_proxy=http://as.lan:1080 --build-arg https_proxy=http://as.lan:1080 -t allape/goview .
# docker tag allape/goview:latest docker-registry.lan.allape.cc/allape/goview:latest
# docker push docker-registry.lan.allape.cc/allape/goview:latest
```

### Run

```shell
docker compose -f docker-compose.yaml up -d
```
