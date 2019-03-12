package main

import (
	"2019_1_Auteam/storage"
	"fmt"
	"net/http"
	"time"
)

const (
	sessionServerAddr = "localhost:8081"
	key               = "server.crt"
	maxUploadSize     = 2 * 1024
)

func main() {
	st, err := storage.OpenPostgreStorage("localhost", "docker", "docker", "docker")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pbClient, err := ConnectToSessionService()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	server := Server{st, pbClient}

	r := CreateRouter(&server)
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	fmt.Println("Start server")
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}
