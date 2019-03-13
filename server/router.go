package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func CreateRouter(server *Server) *mux.Router {
	mux := mux.NewRouter()
	mux.Use(SetCors)
	mux.Use(Logging)
	media := mux.PathPrefix("/media").Subrouter()
	media.HandleFunc("/pic/{id}", server.handleMedia)
	media.HandleFunc("/upload", server.handleUpload)
	userRouter := mux.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", server.handleLogin)
	userRouter.HandleFunc("/signup", server.handleSignup)
	userRouter.HandleFunc("/session", server.handleSession)
	userRouter.HandleFunc("/logout", server.handleLoguot)
	userRouter.HandleFunc("/list", server.handleList)
	userRouter.HandleFunc("/{username}", server.handleUsername)
	userRouter.Handle("/update", server.AuthRequired(http.HandlerFunc(server.handleUserUpdate)))
	return mux
}
