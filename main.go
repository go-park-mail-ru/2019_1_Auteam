package main

import (
    "fmt"
    "strconv"
    "net/http"
    "google.golang.org/grpc"
    "time"
    "encoding/json"
    "crypto/md5"
    "2019_1_Auteam/storage"
    "2019_1_Auteam/models"
    "github.com/gorilla/mux"
)

type Server struct {
    st *storage.PostgreStorage
    grpcCon *grpc.ClientConn
}


func hashPassword(password string) string {
    salt := "lewkrjhnljkdfsgbhkfgjdscwf"
    hash1 := fmt.Sprintf("%x", md5.Sum([]byte(salt + password)))
    hash2 := fmt.Sprint("%x", md5.Sum([]byte(hash1)))
    return fmt.Sprint("%x", md5.Sum([]byte(hash2)))
}

func (s* Server) handleSession(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var request models.LoginRequestJSON
    err := decoder.Decode(&request)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }

    user, err := s.st.GetUserByName(request.Username)
    if err != nil {
        // TODO: send invalid username
    }

    if hashPassword(request.Password) != user.Password {
        // TODO: send invalid password
    }
    // sessionID := s.grpcCon.CreateID(user.ID)
    // expiration := time.Now().Add(2 * 24 * time.Hour)
    // cookie := http.Cookie{Name: "SessionID", Value: sessionID, Expires: expiration}
    // http.SetCoockie(w, &cookie)
    // w.WriteHeader()
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
    var request models.SignUpRequestJSON
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&request)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    user := models.User {Username: request.UserInfo.Username, Email: request.UserInfo.Email, Password: hashPassword(request.Password),}
    s.st.AddUser(&user)
    // expiration := time.Now().Add(2 * 24 * time.Hour)
    // cookie := &http.Cookie{Name: "SessionID", Value: sessionID, Expires: expiration}
    // http.SetCoockie(w, &cookie)
    // w.WriteHeader()
}


func (s *Server) handleLoguot(w http.ResponseWriter, r *http.Request) {
    // cookie, err := r.Cookie("SessionID")
    // if err != nil {
        
    // }
    // w.WriteHeader()
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
    var userPerPage int32 = 10
    page, err := strconv.Atoi(r.FormValue("page"))
    users, err := s.st.GetSortedUsers(int32(page) * userPerPage, userPerPage)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    usersJson := make([]models.AllInfoJSON, 0, 0)
    for _, u := range users {
        userInfo := models.UserInfoJSON{Username: u.Username, Email: u.Email, Userpic: u.Pic}
        gameInfo := models.GameInfoJSON{Score: u.Score}
        info := models.AllInfoJSON{UserInfo: userInfo, GameInfo: gameInfo}
        usersJson = append(usersJson, info)
    }
    err = json.NewEncoder(w).Encode(usersJson)
    w.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) handleUsername(w http.ResponseWriter, r *http.Request) {
    // w.WriteHeader()
}

func main() {
    server := Server{}
    mux := mux.NewRouter()
    s := mux.PathPrefix("/user").Subrouter()
    s.HandleFunc("/login", server.handleLogin).Methods("POST")
    s.HandleFunc("/signup", server.handleSignup).Methods("POST")
    s.HandleFunc("/session", server.handleSession).Methods("GET")
    s.HandleFunc("/logout", server.handleLoguot).Methods("POST")
    s.HandleFunc("/list", server.handleList).Methods("POST")
    s.HandleFunc("/{username}", server.handleUsername).Methods("GET", "POST")
    srv := &http.Server {
        Addr:         "0.0.0.0:8080",
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: mux,
    }
    srv.ListenAndServe()
}