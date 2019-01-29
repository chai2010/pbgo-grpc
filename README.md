# pbgo-examples

- https://github.com/chai2010/pbgo

## Proto File(hello.proto)

```protobuf
syntax = "proto3";

package hello_pb;

import "github.com/chai2010/pbgo/pbgo.proto";

message String {
	string value = 1;
}
message StaticFile {
	string content_type = 1;
	bytes content_body = 2;
}

message Message {
	string value = 1;
	repeated int32 array = 2;
	map<string,string> dict = 3;
	String subfiled = 4;
}

service HelloService {
	option (pbgo.service_opt) = {
		rename: "PBGOHelloService"
	};

	rpc Hello (String) returns (String) {
		option (pbgo.rest_api) = {
			get: "/hello/:value"
			post: "/hello"

			additional_bindings {
				method: "DELETE"; url: "/hello"
			}
			additional_bindings {
				method: "PATCH"; url: "/hello"
			}
		};
	}
	rpc Echo (Message) returns (Message) {
		option (pbgo.rest_api) = {
			get: "/echo/:subfiled.value"
		};
	}
	rpc Static(String) returns (StaticFile) {
		option (pbgo.rest_api) = {
			additional_bindings {
				method: "GET"
				url: "/static/:value"
				content_type: ":content_type"
				content_body: ":content_body"
			}
		};
	}

	rpc ServerStream(String) returns (stream String);
	rpc ClientStream(stream String) returns (String);
	rpc Channel(stream String) returns (stream String);
}
```

Generate stub code:

```
$ go install github.com/chai2010/pbgo-grpc/protoc-gen-pbgo-grpc
$ protoc -I=. -I=./api/third_party --pbgo-grpc_out=plugins=grpc+pbgo:. hello.proto
```

## Example

```go
package main

import (
	"context"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/chai2010/pbgo-grpc/api"
)

var (
	_ pb.HelloServiceServer = (*HelloService)(nil)
)

var (
	flagRoot = flag.String("root", "./testdata", "set root dir")
)

func main() {
	flag.Parse()

	helloService := NewHelloService(*flagRoot)

	go func() {
		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		pb.RegisterHelloServiceServer(grpcServer, helloService)

		lis, err := net.Listen("tcp", ":3999")
		if err != nil {
			log.Fatal(err)
		}
		grpcServer.Serve(lis)
	}()

	ctx := context.Background()
	router := pb.PBGOHelloServiceGrpcHandler(ctx, helloService, nil)
	log.Fatal(http.ListenAndServe(":8080", someMiddleware(router)))
}

func someMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		defer func() {
			timeElapsed := time.Since(timeStart)
			log.Println(r.Method, r.URL, timeElapsed)
		}()

		next.ServeHTTP(wr, r)
	})
}

type HelloService struct {
	rootdir string
}

func NewHelloService(rootdir string) *HelloService {
	return &HelloService{
		rootdir: rootdir,
	}
}

func (p *HelloService) Hello(ctx context.Context, args *pb.String) (*pb.String, error) {
	reply := &pb.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func (p *HelloService) Echo(ctx context.Context, args *pb.Message) (*pb.Message, error) {
	conn, err := grpc.Dial("localhost:3999", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)
	result, err := client.Hello(context.Background(), &pb.String{Value: "hello"})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	reply := &pb.Message{
		Value: result.GetValue(),
	}

	return reply, nil
}

func (p *HelloService) Static(ctx context.Context, args *pb.String) (*pb.StaticFile, error) {
	data, err := ioutil.ReadFile(p.rootdir + "/" + args.Value)
	if err != nil {
		return nil, err
	}

	reply := new(pb.StaticFile)
	reply.ContentType = mime.TypeByExtension(args.Value)
	reply.ContentBody = data
	return reply, nil
}

func (p *HelloService) ServerStream(*pb.String, pb.HelloService_ServerStreamServer) error {
	return errors.New("todo")
}

func (p *HelloService) ClientStream(pb.HelloService_ClientStreamServer) error {
	return errors.New("todo")
}

func (p *HelloService) Channel(pb.HelloService_ChannelServer) error {
	return errors.New("todo")
}
```

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
