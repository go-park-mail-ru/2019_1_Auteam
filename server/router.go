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
	media.HandleFunc("/pic/{id}", server.handleMedia).Methods("GET")
	media.HandleFunc("/upload", server.handleUpload).Methods("POST")
	userRouter := mux.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", server.handleLogin).Methods("POST")
	userRouter.HandleFunc("/signup", server.handleSignup).Methods("POST")
	userRouter.HandleFunc("/session", server.handleSession).Methods("GET")
	userRouter.HandleFunc("/logout", server.handleLoguot).Methods("POST")
	userRouter.HandleFunc("/list", server.handleList).Methods("GET")
	userRouter.HandleFunc("/{username}", server.handleUsername).Methods("GET")
	userRouter.Handle("/update", server.AuthRequired(http.HandlerFunc(server.handleUserUpdate))).Methods("POST")
	return mux
}
