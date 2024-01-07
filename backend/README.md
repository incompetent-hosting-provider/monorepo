# backend
The Gin GoLang REST Backend for our Service

## Quickstart

To build run:
'''sh
go build ./main.go
'''

To run the resulting executable
'''sh
./main [-debug] [-pretty-logs]
'''

## Swagger/OpenAPI

For swagger use [SWAG](https://github.com/swaggo/swag).

After updating the swagger definition run 
```sh
swag init
```

## Tests

To run tests use 
```sh
go test ./...
```