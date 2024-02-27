# Backend
The Gin GoLang REST Backend for our service.

## Quickstart

To build run:
```sh
go build ./main.go
```

To run the resulting executable
```sh
./main [-debug] [-pretty-logs]
```

## Swagger/OpenAPI

For swagger use [SWAG](https://github.com/swaggo/swag).

After updating the swagger definition run 
```sh
swag init
```

When running into issues with the swag installation not being found run:
```sh
export PATH=$(go env GOPATH)/bin:$PATH
```

## Tests

To run tests use 
```sh
go test ./...
```
