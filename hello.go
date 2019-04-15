// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	pbgo_grpc "github.com/chai2010/pbgo-grpc"
	hello_pb "github.com/chai2010/pbgo-grpc/example/api"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	svc := new(HiService)

	// curl -d '{"value":"123"}' 127.0.0.1:8080/hi
	router.POST("/hi", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reply, err := pbgo_grpc.CallMethod(svc.Hi, context.Background(), string(b))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		jsonReply, err := json.Marshal(reply)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		fmt.Fprint(w, string(jsonReply))
		log.Println(string(jsonReply))
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}

type HiService struct{}

func (p *HiService) Hi(ctx context.Context, req *hello_pb.String) (*hello_pb.String, error) {
	log.Printf("HelloService.Hello: req = %v\n", req)

	reply := &hello_pb.String{Value: "hello:" + req.GetValue()}
	return reply, nil
}
