package main

import (
	pb "2019_1_Auteam/protobuf"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCheckID(t *testing.T) {
	sr := sessionRoute{idStorage: map[string]storageData{}}

	_, err := sr.CheckID(context.Background(), &pb.Session{Id: "123"})
	if err == nil {
		t.Error("Should be an error fo nonexistent id")
	}
}

func TestCreateID(t *testing.T) {
	sr := sessionRoute{idStorage: map[string]storageData{}}

	_, err := sr.CreateID(context.Background(), &pb.UserData{UserID: 123})
	if err != nil {
		t.Error("UUID wasn't successfully created")
	}
}

func TestDeleteID(t *testing.T) {
	sr := sessionRoute{idStorage: map[string]storageData{}}

	_, err := sr.DeleteID(context.Background(), &pb.Session{Id: "123"})
	if err == nil {
		t.Error("Unavailable deletion completed")
	}
}

func TestServerStartFunction(t *testing.T) {
	server := ServerStart()
	time.Sleep(time.Millisecond * 50)
	server.GracefulStop()

	if recover() != nil {
		t.Error("Some errors in server startup")
	}
}

func TestComplex(t *testing.T) {
	sr := sessionRoute{idStorage: map[string]storageData{}}

	res, err := sr.CreateID(context.Background(), &pb.UserData{UserID: 777})
	if err != nil {
		t.Error("ID Creation error")
	}
	sessionid := res.Id

	_, err = sr.CheckID(context.Background(), &pb.Session{Id: sessionid})
	if err != nil {
		t.Error("ID Checking error")
	}

	_, err = sr.DeleteID(context.Background(), &pb.Session{Id: sessionid})
	if err != nil {
		t.Error("ID Deletion error")
	}

	_, err = sr.CheckID(context.Background(), &pb.Session{Id: sessionid})
	if err != nil {
		fmt.Println("User session was successfully deleted")
	} else {
		t.Error("ID Deletion check error")
	}
}
