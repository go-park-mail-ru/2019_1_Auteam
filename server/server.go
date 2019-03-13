package server

import (
	pb "2019_1_Auteam/apps/sessions_app/protobuf"
	"2019_1_Auteam/storage"
	"context"
	"log"
	"fmt"
	"google.golang.org/grpc"
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
	log.Println("GRPC DELETE err = ", err)
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
	log.Println("GRPC CHECK err =", err, "res =", res)
	return res.UserID, err
}

func (s *Server) CreateSession(userId int32) (string, error) {
	if s.SessionClient == nil {
		return "", fmt.Errorf("No connection to grpc")
	}
	res, err := s.SessionClient.CreateID(context.Background(), &pb.UserData{UserID: userId})
	log.Println("GRPC CREATE err =", err, "res =", res)
	if err != nil {
		return "", err
	}
	return res.Id, nil
}

func ConnectToSessionService() (client pb.SessionRouteClient, err error) {
	conn, err := grpc.Dial(sessionServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return
	}
	client = pb.NewSessionRouteClient(conn)
	return
}
