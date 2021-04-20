package main

import (
	"deyforyou/dey/schema"
	"deyforyou/dey/service"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

var server *grpc.Server

func init() {
	server = grpc.NewServer()
}

func init() {
	schema.RegisterArticleServiceServer(server, service.NewArticleServiceServer())
}

func main() {
	listener, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func runAfter(d time.Duration, f func()) {
	t := time.NewTicker(d)
	f()
	for {
		select {
		case <-t.C:
			go f()
		}
	}
}
