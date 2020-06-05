# go-grpc-example

Just a gRPC example.

###### Server
```
go run server/server.go
```

###### Client
```sh
# Put a new entry (returns old value)
go run client/client.go put --key PI --value 3.14
# Get a value
go run client/client.go get --key PI
```