package main

import(
	"net/http"
	"io"
	"os"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"fmt"
	"2019_1_Auteam/models"
	"encoding/json"
)

func (s *Server) handleMedia(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	f, err := os.OpenFile("./media/"+(mux.Vars(r))["id"], os.O_RDONLY, 0666)
	if err != nil {
		w.WriteHeader(404)
	}
	defer f.Close()
	io.Copy(w, f)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upload file")
	r.ParseMultipartForm(maxUploadSize)
	file, handler, err := r.FormFile("my_file")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	defer file.Close()
	mimeType := handler.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg":
	case "image/png":
	default:
		w.WriteHeader(400)
		return
	}
	fileId := uuid.New()
	f, err := os.OpenFile("./media/"+fileId.String(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	userpicJson := models.UserPicJSON{fileId.String()}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(userpicJson)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}