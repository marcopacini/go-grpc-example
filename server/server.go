package main

import (
	"context"
	"go-grpc-example/cache"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	entries map[string]string
}

func newServer(size int) *server {
	return &server{
		entries: make(map[string]string, size),
	}
}

func (s *server) Put(_ context.Context, e *cache.Entry) (*cache.Value, error) {
	old, _ := s.entries[e.Key]
	s.entries[e.Key] = e.Value
	return &cache.Value{Value: old}, status.New(codes.OK, "").Err()
}

func (s *server) Get(_ context.Context, k *cache.Key) (*cache.Value, error) {
	if v, ok := s.entries[k.Key]; ok {
		return &cache.Value{Value: v}, status.New(codes.OK, "").Err()
	}
	return &cache.Value{}, status.Errorf(codes.NotFound, "key does not exist: %v", k.Key)
}

const (
	port      = ":50051"
	sizeCache = 1024
)

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", port)
	}

	s := grpc.NewServer()
	cache.RegisterCacheServer(s, newServer(sizeCache))

	log.Printf("Starting gRPC listener on port %v", port)
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
