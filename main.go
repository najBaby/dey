package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
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
	port := "9000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	fmt.Printf("Starting to run on %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("hello"))
	}))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
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
