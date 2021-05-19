package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"deyforyou/dey/schema"
	"deyforyou/dey/service"

	"google.golang.org/grpc"
)

var server *grpc.Server

func init() {
	server = grpc.NewServer()
}

func init() {
	schema.RegisterMovieServiceServer(server, service.NewMovieServiceServer())
}

func main() {
	log.Printf("Starting to run on %s", "9000")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "9000"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func Run(d time.Duration, f func()) {
	f()
	group := new(sync.WaitGroup)
	for range time.NewTicker(d).C {
		group.Add(1)
		go func() {
			f()
			group.Done()
		}()
	}
	group.Wait()
}
