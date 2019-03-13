package server

import (
	pb "2019_1_Auteam/apps/sessions_app/protobuf"
	"2019_1_Auteam/storage"
	"context"
	"log"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	sessionServerAddr = "sessions_server:8081"
	key               = "server.crt"
	maxUploadSize     = 2 * 1024
)

type Server struct {
	St            storage.StorageI
	SessionClient pb.SessionRouteClient
}

func (s *Server) DeleteSession(session string) error {
	if s.SessionClient == nil {
		return fmt.Errorf("No connection to grpc")
	}
	_, err := s.SessionClient.DeleteID(context.Background(), &pb.Session{Id: session})
	return err
}

func (s *Server) CheckSession(session string) (int32, error) {
	if s.SessionClient == nil {
		return 0, fmt.Errorf("No connection to grpc")
	}
	res, err := s.SessionClient.CheckID(context.Background(), &pb.Session{Id: session})
	if err != nil {
		return 0, err
	}
	return res.UserID, err
}

func (s *Server) CreateSession(userId int32) (string, error) {
	if s.SessionClient == nil {
		return "", fmt.Errorf("No connection to grpc")
	}
	res, err := s.SessionClient.CreateID(context.Background(), &pb.UserData{UserID: 777})
	if err != nil {
		return "", err
	}
	return res.Id, nil
}

func ConnectToSessionService() (client pb.SessionRouteClient, err error) {
	creds, err := credentials.NewClientTLSFromFile(key, "")
	if err != nil {
		log.Println(err.Error())
		return
	}
	conn, err := grpc.Dial(sessionServerAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Println(err.Error())
		return
	}
	client = pb.NewSessionRouteClient(conn)
	return
}
