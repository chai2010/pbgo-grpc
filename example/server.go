// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	ctxpkg "github.com/chai2010/pbgo-grpc/context"
	pb "github.com/chai2010/pbgo-grpc/example/api"
)

var (
	_ pb.HelloServiceServer = (*HelloService)(nil)
)

var (
	flagRoot = flag.String("root", "./testdata", "set root dir")
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	flag.Parse()

	helloService := NewHelloService(*flagRoot)

	go func() {
		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		pb.RegisterHelloServiceServer(grpcServer, helloService)

		log.Println("grpc server on :3999")
		lis, err := net.Listen("tcp", ":3999")
		if err != nil {
			log.Fatal(err)
		}
		grpcServer.Serve(lis)
	}()

	ctx := context.Background()
	router := pb.PBGOHelloServiceGrpcHandler(
		ctx, helloService,
		func(ctx context.Context, req *http.Request) (context.Context, error) {
			return ctxpkg.AnnotateContext(ctx, req, nil)
		},
	)

	log.Println("http server on :8080")
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

func (p *HelloService) Hello(ctx context.Context, req *pb.String) (*pb.String, error) {
	log.Printf("HelloService.Hello: req = %v\n", req)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md: %v\n", md)
	} else {
		log.Println("no metadata")
	}

	reply := &pb.String{Value: "hello:" + req.GetValue()}
	return reply, nil
}

func (p *HelloService) Echo(ctx context.Context, req *pb.Message) (*pb.Message, error) {
	log.Printf("HelloService.Echo: req = %v\n", req)

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

func (p *HelloService) Static(ctx context.Context, req *pb.String) (*pb.StaticFile, error) {
	log.Printf("HelloService.Static: req = %v\n", req)

	data, err := ioutil.ReadFile(p.rootdir + "/" + req.Value)
	if err != nil {
		return nil, err
	}

	reply := new(pb.StaticFile)
	reply.ContentType = mime.TypeByExtension(req.Value)
	reply.ContentBody = data
	return reply, nil
}

func (p *HelloService) ServerStream(*pb.String, pb.HelloService_ServerStreamServer) error {
	log.Printf("HelloService.ServerStream: todo\n")

	return errors.New("todo")
}

func (p *HelloService) ClientStream(pb.HelloService_ClientStreamServer) error {
	log.Printf("HelloService.ClientStream: todo\n")

	return errors.New("todo")
}

func (p *HelloService) Channel(pb.HelloService_ChannelServer) error {
	log.Printf("HelloService.Channel: todo\n")

	return errors.New("todo")
}
