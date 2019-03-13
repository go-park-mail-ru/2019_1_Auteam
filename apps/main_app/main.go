package main

import (
	"2019_1_Auteam/server"
	"2019_1_Auteam/storage"
	"fmt"
	"net/http"
	"time"
	"os"
	"log"
)

func main() {
	log.SetOutput(os.Stdout)

	st, err := storage.OpenPostgreStorage("database", "docker", "docker", "docker")
	if err != nil {
		log.Println("Can't connect to database")
		log.Println(err.Error())
		return
	}

	pbClient, err := server.ConnectToSessionService()
	if err != nil {
		log.Println("Can't connect to session service")
		log.Println(err.Error())
		pbClient = nil
	}

	s := server.Server{st, pbClient}

	r := server.CreateRouter(&s)
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
		log.Println(err.Error())
	}
}
