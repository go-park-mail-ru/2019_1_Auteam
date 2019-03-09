package main

import (
    pb "./proto"
    "flag"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "log"
    "net"
    "os"
)

var (
    crt = "../server.crt"
    key = "../server.key"
)

func main() {
    log.SetOutput(os.Stdout)
    flag.Parse()
    lis, err := net.Listen("tcp", ":8081")
    if err != nil {
        log.Println("Cannot bind port")
        panic(err)
    }

    creds, err := credentials.NewServerTLSFromFile(crt, key)
    if err != nil {
        log.Println("Cannot load TLS configuration")
        panic(err)
    }
    
    server := grpc.NewServer(grpc.Creds(creds))
    pb.RegisterSessionRouteServer(server, &SessionRoute{idStorage: map[string]storageData{}})
    err = server.Serve(lis)
    if err != nil {
        log.Println("Cannot start server")
        panic(err)
    }
}
