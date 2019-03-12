package main

import (
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
	server, err := NewServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	r := CreateRouter(server)
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	fmt.Println("Start server")
	srv.ListenAndServe()
}
