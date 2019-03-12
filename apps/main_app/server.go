package main

import (
	pb "2019_1_Auteam/apps/sessions_app/protobuf"
	"2019_1_Auteam/storage"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	st            storage.StorageI
	sessionClient pb.SessionRouteClient
}

func (s *Server) DeleteSession(session string) error {
	_, err := s.sessionClient.DeleteID(context.Background(), &pb.Session{Id: session})
	return err
}

func (s *Server) CheckSession(session string) (int32, error) {
	res, err := s.sessionClient.CheckID(context.Background(), &pb.Session{Id: session})
	if err != nil {
		return 0, err
	}
	return res.UserID, err
}

func (s *Server) CreateSession(userId int32) (string, error) {
	res, err := s.sessionClient.CreateID(context.Background(), &pb.UserData{UserID: 777})
	if err != nil {
		return "", err
	}
	return res.Id, nil
}

func ConnectToSessionService() (client pb.SessionRouteClient, err error) {
	creds, err := credentials.NewClientTLSFromFile(key, "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn, err := grpc.Dial(sessionServerAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client = pb.NewSessionRouteClient(conn)
	return
}
