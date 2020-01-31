package main

import (
	"cache/service/cache"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)


type Method string

const (
	Get Method = "Get"
	Put Method = "Put"
)

var serv *grpc.Server

const (
	address = ":50051"
)

func tearUp() error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", address)
	}

	serv = grpc.NewServer()
	cache.RegisterCacheServer(serv, newServer(sizeCache))

	go func () {
		if err := serv.Serve(l); err != nil {
			fmt.Printf("failed to serve: %v", err)
		}
	}()

	return nil
}

func tearDown() {
	serv.Stop()
}

func TestServer(t *testing.T) {
	tcs := []struct {
		key string
		value string
		method Method
	} {
		{ key: "batman", value: "Bruce Wayne", method: Put },
		{ key: "superman", value: "Clark Kent", method: Put },
		// { key: "batman", value: "Bruce Wayne", method: Get },
	}

	if err := tearUp(); err != nil {
		t.Errorf("server not started: %v", err)
	}
	defer tearDown()

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Errorf("not connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Errorf("not close: %v", err)
		}
	}()

	client := cache.NewCacheClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, tc := range tcs {
		e := &cache.Entry{
			Key:	tc.key,
			Value:  tc.value,
		}

		v, err := client.Put(ctx, e)
		if err != nil {
			t.Error(err)
		}

		if tc.method == Get && tc.value != v.Value {
			t.Errorf("get %v, want %v", v.Value, tc.value)
		}
	}
}
