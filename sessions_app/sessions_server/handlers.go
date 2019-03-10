package main

import (
	pb "../protobuf"
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
)

type storageData struct {
	id int32
}

type sessionRoute struct {
	idStorage map[string]storageData
}

func (sr *sessionRoute) CreateID(ctx context.Context, in *pb.UserData) (*pb.Session, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println("UUID generation error")
		panic(err)
	}

	sr.idStorage[id.String()] = storageData{id: in.UserID}
	return &pb.Session{Id: id.String()}, nil
}

func (sr *sessionRoute) CheckID(ctx context.Context, in *pb.Session) (*pb.UserData, error) {
	data, exists := sr.idStorage[in.Id]
	if !exists {
		return nil, errors.New("session-id not exists")
	}

	return &pb.UserData{UserID: data.id}, nil
}

func (sr *sessionRoute) DeleteID(ctx context.Context, in *pb.Session) (*pb.Empty, error) {
	_, exists := sr.idStorage[in.Id]
	if !exists {
		return nil, errors.New("session-id not exists")
	}

	delete(sr.idStorage, in.Id)
	return &pb.Empty{}, nil
}
