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
	media.HandleFunc("/image/", server.handleUpload)
	media.HandleFunc("/image/{id}", server.handleMedia)
	sessionRouter := mux.PathPrefix("/session").Subrouter()
	sessionRouter.HandleFunc("", server.handleLogin).Methods("POST")
	sessionRouter.Handle("", server.AuthRequired(http.HandlerFunc(server.handleUserUpdate))).Methods("PATCH")
	sessionRouter.HandleFunc("", server.handleLoguot).Methods("DELETE")
	sessionRouter.HandleFunc("", server.handleSession).Methods("GET")

	userRouter := mux.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", server.handleList).Methods("GET")
	userRouter.HandleFunc("", server.handleSignup).Methods("PUT")
	userRouter.HandleFunc("/{username}", server.handleUsername).Methods("GET")
	return mux
}
