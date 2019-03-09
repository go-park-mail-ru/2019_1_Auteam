package main

import (
    "fmt"
    "context"
    "google.golang.org/grpc"
    pb "../sessions/proto"
    "google.golang.org/grpc/credentials"
)

var (
    key = "server.crt"
    adr = "localhost:8081"
)

func main() {
    creds, err := credentials.NewClientTLSFromFile(key, "")
    if err != nil {
        panic(err)
    }

    fmt.Println(creds)

    conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(creds))
    if err != nil {
        fmt.Println("gRPC server is not available")
        panic(err)
    }

    client := pb.NewSessionRouteClient(conn)

    res, err := client.CreateID(context.Background(), &pb.UserData{UserID: 777})
    if err != nil {
        panic(err)
    }
    fmt.Println("Check create ID: ", res.Id)

    res2, err := client.CheckID(context.Background(), &pb.Session{Id: res.Id})
    if err != nil {
        panic(err)
    }

    fmt.Println("Get id", res2.UserID)

    defer conn.Close()
}
