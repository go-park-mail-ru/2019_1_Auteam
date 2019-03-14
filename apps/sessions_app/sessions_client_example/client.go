package main

import (
	pb "2019_1_Auteam/apps/sessions_app/protobuf"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
)

var (
	adr = "localhost:8081"
)

func main() {
	log.SetOutput(os.Stdout)

	//Connect to server
	conn, err := grpc.Dial(adr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("gRPC server is not available")
		panic(err)
	}
	client := pb.NewSessionRouteClient(conn)

	//Create ID example
	res, err := client.CreateID(context.Background(), &pb.UserData{UserID: 777})
	if err != nil {
		panic(err)
	}
	sessionid := res.Id
	fmt.Println("Check create ID: ", res.Id)

	//Check ID example
	res2, err := client.CheckID(context.Background(), &pb.Session{Id: sessionid})
	if err != nil {
		panic(err)
	}
	fmt.Println("Get id", res2.UserID)

	//Delete ID example
	_, err = client.DeleteID(context.Background(), &pb.Session{Id: sessionid})

	res2, err = client.CheckID(context.Background(), &pb.Session{Id: sessionid})
	if err != nil {
		fmt.Println("User session was successfully deleted")
	} else {
		fmt.Println("User were not successfully deleted. Id: ", res2.UserID)
	}

	defer conn.Close()
}
