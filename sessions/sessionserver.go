package main

import (
    "context"
    "fmt"
    pb "./proto"
    "github.com/google/uuid"
    "errors"
)

type SessionRoute struct {
}

type StorageUserData struct {
    id int32
}

var SessionIDStorage = map[string]StorageUserData{}

func (sr *SessionRoute) CreateID(ctx context.Context, in *pb.UserData) (*pb.Session, error) {
    id, err := uuid.NewUUID()
    if (err != nil) {
        panic(err)
    }

    SessionIDStorage[id.String()] = StorageUserData{id: in.UserID}

    fmt.Println("CreateID: ", in.UserID)
    fmt.Println("Created UUID: ", id.String())
    return &pb.Session{Id: id.String()}, nil
}

func (sr *SessionRoute) CheckID(ctx context.Context, in *pb.Session) (*pb.UserData, error) {
    data, exists := SessionIDStorage[in.Id]
    if !exists {
        return nil, errors.New("session-id not exists")
    }

    fmt.Println("CheckID ", in.Id)
    return &pb.UserData{UserID: data.id}, nil
}

func (sr *SessionRoute) DeleteID(ctx context.Context, in *pb.Session) (*pb.Empty, error) {
    fmt.Println("DeleteID", in.Id)
    _, exists := SessionIDStorage[in.Id]
    if !exists {
        return nil, errors.New("session-id not exists")
    }

    delete(SessionIDStorage, in.Id)
    return &pb.Empty{}, nil
}
