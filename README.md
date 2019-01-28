# pbgo-examples

Start gRPC & Rest server:

```
go run server.go
```

Rest API:

```
$ curl localhost:8080/hello/gopher
{"value":"hello:gopher"}
$ curl localhost:8080/hello/gopher?value=vgo
{"value":"hello:vgo"}
$ curl localhost:8080/hello -X POST --data '{"value":"cgo"}'
{"value":"hello:cgo"}

$ curl localhost:8080/echo/gopher
{"subfiled":{"value":"gopher"}}
$ curl "localhost:8080/echo/gopher?array=123&array=456"
{"array":[123,456],"subfiled":{"value":"gopher"}}
$ curl "localhost:8080/echo/gopher?dict%5Babc%5D=123"
{"dict":{"abc":"123"},"subfiled":{"value":"gopher"}}

$ curl localhost:8080/static/gopher.png
$ curl localhost:8080/static/hello.txt
```

gRPC API:

```
$ grpcurl -plaintext localhost:3999 list
```

## BUGS

Report bugs to <chaishushan@gmail.com>.

Thanks!
