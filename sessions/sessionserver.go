package main

import (
	pb "./proto"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type storageData struct {
	id int32
}

type SessionRoute struct {
	idStorage map[string]storageData
}

func (sr *SessionRoute) CreateID(ctx context.Context, in *pb.UserData) (*pb.Session, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println("UUID generation error")
		panic(err)
	}

	sr.idStorage[id.String()] = storageData{id: in.UserID}
	return &pb.Session{Id: id.String()}, nil
}

func (sr *SessionRoute) CheckID(ctx context.Context, in *pb.Session) (*pb.UserData, error) {
	data, exists := sr.idStorage[in.Id]
	if !exists {
		return nil, errors.New("session-id not exists")
	}

	return &pb.UserData{UserID: data.id}, nil
}

func (sr *SessionRoute) DeleteID(ctx context.Context, in *pb.Session) (*pb.Empty, error) {
	fmt.Println("DeleteID", in.Id)
	_, exists := sr.idStorage[in.Id]
	if !exists {
		return nil, errors.New("session-id not exists")
	}

	delete(sr.idStorage, in.Id)
	return &pb.Empty{}, nil
}
