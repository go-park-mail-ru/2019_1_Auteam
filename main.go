package main

import (
    "fmt"
    "os"
    "io"
    "strconv"
    "net/http"
    "google.golang.org/grpc"
    "time"
    "encoding/json"
    "crypto/md5"
    "2019_1_Auteam/storage"
    "google.golang.org/grpc/credentials"
    "2019_1_Auteam/models"
    "github.com/gorilla/mux"
    "context"
    "github.com/google/uuid"
    pb "2019_1_Auteam/sessions_app/protobuf"
)

const (
    sessionServerAddr = "localhost:8081"
    key = "server.crt"
    maxUploadSize = 2 * 1024
)

type Server struct {
    st *storage.PostgreStorage
    sessionClient pb.SessionRouteClient
}

func (s *Server) DeleteSession(session string) error {
    _, err := s.sessionClient.DeleteID(context.Background(), &pb.Session{Id: session})
    return err
}

func (s *Server) CheckSession(session string) (int32, error) {
    res, err := s.sessionClient.CheckID(context.Background(), &pb.Session{Id: session})
    if err != nil {
        return 0, err
    }
    return res.UserID, err
}

func (s *Server) CreateSession(userId int32) (string, error) {
    res, err := s.sessionClient.CreateID(context.Background(), &pb.UserData{UserID: 777})
    if err != nil {
        return "", err
    }
    return res.Id, nil
}

func NewServer() (*Server, error) {
    creds, err := credentials.NewClientTLSFromFile(key, "")
    if err != nil {
        return nil, err
    }
    conn, err := grpc.Dial(sessionServerAddr, grpc.WithTransportCredentials(creds))
    if err != nil {
        return nil, err
    }
    client := pb.NewSessionRouteClient(conn)
    st, err := storage.OpenPostgreStorage("postgres://e.kislukha@localhost/back_db?sslmode=disable")
    if err != nil {
        return nil, err
    }
    return &Server{st, client}, nil
}

func hashPassword(password string) string {
    salt := "lewkrjhnljkdfsgbhkfgjdscwf"
    hash1 := fmt.Sprintf("%x", md5.Sum([]byte(salt + password)))
    hash2 := fmt.Sprint("%x", md5.Sum([]byte(hash1)))
    return fmt.Sprint("%x", md5.Sum([]byte(hash2)))
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("SessionID")
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    userId, err := s.CheckSession(cookie.Value)
    user, err := s.st.GetUserById(userId)
    if err != nil {
        w.WriteHeader(400)
        return
    }
    userInfo := models.UserInfoJSON{Username: &user.Username, Email: &user.Email, Userpic: &user.Pic}
    gameInfo := models.GameInfoJSON{Score: &user.Score}
    info := models.AllInfoJSON{UserInfo: &userInfo, GameInfo: &gameInfo}
    encoder := json.NewEncoder(w)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    err = encoder.Encode(info)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var request models.LoginRequestJSON
    err := decoder.Decode(&request)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }

    user, err := s.st.GetUserByName(*request.Username)
    if err != nil {
        w.WriteHeader(400)
        return
    }

    if hashPassword(*request.Password) != user.Password {
        w.WriteHeader(400)
        return
    }

    res, err := s.CreateSession(user.ID)
    if err != nil {
        w.WriteHeader(500)
    }
    expiration := time.Now().Add(2 * 24 * time.Hour)
    cookie := http.Cookie{Name: "SessionID", Value: res, Expires: expiration}
    http.SetCookie(w, &cookie)
    w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
    var request models.SignUpRequestJSON
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&request)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    
    user := models.User{Username: *request.UserInfo.Username, Email: *request.UserInfo.Email, Password: hashPassword(*request.Password),}
    err = s.st.AddUser(&user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)   
    }

    sessionId, err := s.CreateSession(user.ID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    expiration := time.Now().Add(2 * 24 * time.Hour)
    cookie := http.Cookie{Name: "SessionID", Value: sessionId, Expires: expiration}
    http.SetCookie(w, &cookie)
    response := models.SignUpResponseJSON{
        &models.ValidateJSON{
            Success: true,
        },
        &models.ValidateJSON{
            Success: true,
        },
        &models.ValidateJSON{
            Success: true,
        },
        &models.ValidateJSON{
            Success: true,
        },
        &models.ErrorJSON{},
    }

    encoder := json.NewEncoder(w)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    err = encoder.Encode(response)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    w.WriteHeader(http.StatusOK)
}


func (s *Server) handleLoguot(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("SessionID")
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
    }
    c := http.Cookie {
        Name: "SessionID",
        MaxAge: 0,
        Value: "",
    }
    http.SetCookie(w, &c)
    sessionId := cookie.Value
    s.DeleteSession(sessionId)
    w.WriteHeader(http.StatusOK)
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
        userInfo := models.UserInfoJSON{Username: &u.Username, Email: &u.Email, Userpic: &u.Pic}
        gameInfo := models.GameInfoJSON{Score: &u.Score}
        info := models.AllInfoJSON{UserInfo: &userInfo, GameInfo: &gameInfo}
        usersJson = append(usersJson, info)
    }
    err = json.NewEncoder(w).Encode(usersJson)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
    var request models.UpdateUserRequestJSON
    response := models.UpdateResponseJSON{
        &models.ValidateJSON {
            Success: true,
        },
        &models.ValidateJSON {
            Success: true,
        },
        &models.ValidateJSON {
            Success: true,
        },
        &models.ValidateJSON {
            Success: true,
        },
        &models.ValidateJSON {
            Success: true,
        },
        &models.ErrorJSON {
        },
    }
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&request)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    userId := r.Context().Value("userID").(int32)
    user, err := s.st.GetUserById(userId)
    if err != nil {
        w.WriteHeader(404)
        return
    }
    if request.OldPass != nil {
        if hashPassword(*request.OldPass) != user.Password {
            response.OldPassValidate.Success = false
        } else {
            if request.NewPass != nil {
                err := s.st.ChangePassword(userId, hashPassword(*request.NewPass))
                if err != nil {
                    w.WriteHeader(http.StatusInternalServerError)
                    return
                }
            }
        }
    }
    if request.UserInfo.Username != nil {
        err := s.st.ChangeUsername(userId, *request.UserInfo.Username)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }
    if request.UserInfo.Userpic != nil {
        err := s.st.ChangePic(userId, *request.UserInfo.Userpic)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }
    if request.UserInfo.Email != nil {
        err := s.st.ChangeEmail(userId, *request.UserInfo.Username)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }
    encoder := json.NewEncoder(w)
    err = encoder.Encode(response)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUsername(w http.ResponseWriter, r *http.Request) {
    username := (mux.Vars(r))["username"]
    user, err := s.st.GetUserByName(username)
    if err != nil {
        w.WriteHeader(404)
        return
    }
    userJson := models.AllInfoJSON{
        &models.UserInfoJSON{
            Username: &user.Username,
            Userpic: &user.Pic,
            Email: &user.Email,
        },
        &models.GameInfoJSON{
            Score: &user.Score,
        },
    }
    encoder := json.NewEncoder(w)
    err = encoder.Encode(userJson)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    w.WriteHeader(http.StatusOK)
}

func (s *Server) handleMedia(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
    f, err := os.OpenFile("./media/" + (mux.Vars(r))["id"], os.O_RDONLY, 0666)
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
    f, err := os.OpenFile("./media/" + fileId.String(), os.O_WRONLY | os.O_CREATE, 0666)
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

func main() {
    server, err := NewServer()
    if err != nil {
        fmt.Println(err)
        return
    }
    mux := mux.NewRouter()
    media := mux.PathPrefix("/media").Subrouter()
    media.HandleFunc("/pic/{id}", server.handleMedia).Methods("GET")
    media.HandleFunc("/upload", server.handleUpload).Methods("POST")
    userRouter := mux.PathPrefix("/user").Subrouter()
    userRouter.HandleFunc("/login", server.handleLogin).Methods("POST")
    userRouter.HandleFunc("/signup", server.handleSignup).Methods("POST")
    userRouter.HandleFunc("/session", server.handleSession).Methods("GET")
    userRouter.HandleFunc("/logout", server.handleLoguot).Methods("POST")
    userRouter.HandleFunc("/list", server.handleList).Methods("POST")
    userRouter.HandleFunc("/{username}", server.handleUsername).Methods("GET")
    userRouter.Handle("/update", server.AuthRequired(http.HandlerFunc(server.handleUserUpdate))).Methods("POST")

    srv := &http.Server{
        Addr:         "0.0.0.0:8080",
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: mux,
    }
    fmt.Println("Start server")
    srv.ListenAndServe()
}