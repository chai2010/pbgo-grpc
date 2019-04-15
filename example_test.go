// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pbgo_grpc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pbgo_grpc "github.com/chai2010/pbgo-grpc"
	hello_pb "github.com/chai2010/pbgo-grpc/example/api"
)

func ExampleCallMethod() {
	fn := func(ctx context.Context, req *hello_pb.String) (*hello_pb.String, error) {
		reply := &hello_pb.String{Value: "hello " + req.GetValue()}
		return reply, nil
	}

	reply, err := pbgo_grpc.CallMethod(fn, context.Background(), `{"value":"9527"}`)
	if err != nil {
		log.Fatal(err)
	}

	jsonReply, err := json.Marshal(reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%T\n", reply)
	fmt.Printf("%s\n", reply.(*hello_pb.String).Value)
	fmt.Printf("%s\n", jsonReply)

	// Output:
	// *hello_pb.String
	// hello 9527
	// {"value":"hello 9527"}
}
