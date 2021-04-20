package main

import (
	"deyforyou/dey/schema"
	"deyforyou/dey/service"
	"log"
	"net"
	"os"
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
	log.Printf("Starting to run on %s", os.Getenv("PORT"))
	listener, err := net.Listen("tcp", os.Getenv("PORT"))
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
