# Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

GOPATH:=$(shell go env GOPATH)

default:
	mkdir -p ./third_party/github.com/chai2010/pbgo
	go run ./tools/sync_proto_file.go > ./third_party/github.com/chai2010/pbgo/pbgo.proto

	go fmt github.com/chai2010/pbgo-examples/protoc-gen-pbgo-examples/...
	go install github.com/chai2010/pbgo-examples/protoc-gen-pbgo-examples/...
	protoc -I=. -I=./third_party --pbgo-examples_out=plugins=grpc+pbgo:. hello.proto
	go fmt && go vet && go test

protoc-gen-pbgo:
	go install github.com/chai2010/pbgo-examples/protoc-gen-pbgo-examples

clean:
	-rm *.pb.go
