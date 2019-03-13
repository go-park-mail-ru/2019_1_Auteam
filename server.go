package main

import (
	pb "2019_1_Auteam/protobuf"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	network = "tcp"
	address = "localhost:8081"
	crt     = "server.crt"
	key     = "server.key"
)

func ServerStart() *grpc.Server {
	lis, err := net.Listen(network, address)
	if err != nil {
		log.Println("Cannot bind port")
		panic(err)
	}
	server := grpc.NewServer()
	pb.RegisterSessionRouteServer(server, &sessionRoute{idStorage: map[string]storageData{}})

	go func() {
		err = server.Serve(lis)
		if err != nil {
			log.Println("Cannot start server")
			panic(err)
		}
	}()

	fmt.Println("Server successfully started")
	return server
}

func main() {
	log.SetOutput(os.Stderr)
	ServerStart()

	// future architecture will support commands like:
	// server start
	// server stop
	// ... etc
	quit := make(chan struct{})
	<-quit
}
