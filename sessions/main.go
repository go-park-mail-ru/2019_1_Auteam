package main

import (
    "flag"
    "net"
    "fmt"
    pb "./proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

var (
    crt = "server.crt"
    key = "server.key"
)

func main() {
    flag.Parse()
    lis, err := net.Listen("tcp", ":8081")
    if err != nil {
        fmt.Println("Cannot bind port. Fatal error")
        panic(err)
    }

    creds, err := credentials.NewServerTLSFromFile(crt, key)
    if err != nil {
        fmt.Errorf("could not load TLS keys: %s", err)
        panic(err)
    }
    
    server := grpc.NewServer(grpc.Creds(creds))

    pb.RegisterSessionRouteServer(server, &SessionRoute{})
    server.Serve(lis)
}
