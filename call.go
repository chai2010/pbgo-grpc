// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pbgo_grpc

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
)

func CallMethod(fn interface{}, ctx context.Context, jsonRequest string) (proto.Message, error) {
	return CallMethodEx(fn, ctx, func(req proto.Message) error {
		return json.Unmarshal([]byte(jsonRequest), req)
	})
}

func CallMethodEx(fn interface{}, ctx context.Context, decodeRequest func(req proto.Message) error) (proto.Message, error) {
	fnType := reflect.TypeOf(fn)

	// check fn type
	if err := checkGrpcMethod(fnType); err != nil {
		return nil, err
	}

	// parse request
	request := reflect.New(fnType.In(1).Elem()).Interface().(proto.Message)
	if err := decodeRequest(request); err != nil {
		return nil, err
	}

	// call grpc method
	// response, err := func(ctx, request)
	returns := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(request),
	})
	if len(returns) != 2 {
		return nil, errors.New("pbfo_grpc: fn must return 2 values")
	}
	if errx := returns[1].Interface(); errx != nil {
		return nil, errx.(error)
	}

	// OK
	response := returns[0].Interface().(proto.Message)
	return response, nil
}

func checkGrpcMethod(fnType reflect.Type) error {
	if fnType.Kind() != reflect.Func {
		return errors.New("pbgo_grpc: fn must be func type")
	}

	// func(ctx context.Context, request proto.Message) (proto.Message, error)
	if fnType.NumIn() != 2 {
		return errors.New("pbgo_grpc: fn must have 2 intput arguments")
	}
	if fnType.NumOut() != 2 {
		return errors.New("pbgo_grpc: fn must return 2 values")
	}

	// in: ctx context.Context
	// if _, ok := fnType.In(0).(context.Context); ok { ... }
	if ok := fnType.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()); !ok {
		return errors.New("pbgo_grpc: fn argument[0] must be `context.Context` type")
	}

	// in: request proto.Message
	// if _, ok := fnType.In(1).(proto.Message); ok { ... }
	if ok := fnType.In(1).Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()); !ok {
		return errors.New("pbgo_grpc: fn argument[1] must be `proto.Message` type")
	}

	// return[0] proto.Message
	// if _, ok := fnType.Out(0).(proto.Message); ok { ... }
	if ok := fnType.Out(0).Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()); !ok {
		return errors.New("pbgo_grpc: fn return[0] must be `proto.Message` type")
	}

	// return[1] error
	// if _, ok := fnType.Out(1).(error); ok { ... }
	if ok := fnType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()); !ok {
		return errors.New("pbgo_grpc: fn return[1] must be `error` type")
	}

	return nil
}
