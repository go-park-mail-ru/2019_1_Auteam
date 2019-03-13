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
	media.HandleFunc("/pic/{id}", server.handleMedia).Methods("GET, OPTIONS")
	media.HandleFunc("/upload", server.handleUpload).Methods("POST, OPTIONS")
	userRouter := mux.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", server.handleLogin).Methods("POST, OPTIONS")
	userRouter.HandleFunc("/signup", server.handleSignup).Methods("POST, OPTIONS")
	userRouter.HandleFunc("/session", server.handleSession).Methods("GET, OPTIONS")
	userRouter.HandleFunc("/logout", server.handleLoguot).Methods("POST, OPTIONS")
	userRouter.HandleFunc("/list", server.handleList).Methods("GET, OPTIONS")
	userRouter.HandleFunc("/{username}", server.handleUsername).Methods("GET, OPTIONS")
	userRouter.Handle("/update", server.AuthRequired(http.HandlerFunc(server.handleUserUpdate))).Methods("POST, OPTIONS")
	return mux
}
