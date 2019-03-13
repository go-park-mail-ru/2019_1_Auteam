package server

import (
	"2019_1_Auteam/models"
	"2019_1_Auteam/validation_tools"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
	"log"
)

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId, err := s.CheckSession(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := s.St.GetUserById(userId)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	userInfo := models.UserInfoJSON{Username: &user.Username, Email: &user.Email, Userpic: &user.Pic}
	score := strconv.Itoa(int(user.Score))
	gameInfo := models.GameInfoJSON{Score: &score}
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

	user, err := s.St.GetUserByName(*request.Username)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if HashPassword(*request.Password) != user.Password {
		w.WriteHeader(400)
		return
	}

	res, err := s.CreateSession(user.ID)
	if err != nil {
		log.Println(err.Error())
	}
	expiration := time.Now().Add(2 * 24 * time.Hour)
	cookie := http.Cookie{Name: "SessionID", Value: res, Expires: expiration}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	isValidRequest := true
	
	var request models.SignUpRequestJSON
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	response := models.SignUpResponseJSON{
		UsernameValidate: &models.ValidateJSON{
			Success: true,
		},
		EmailValidate: &models.ValidateJSON{
			Success: true,
		},
		PasswordValidate: &models.ValidateJSON{
			Success: true,
		},
	}

	if request.UserInfo.Username == nil || *(request.UserInfo.Username) == "" {
		response.UsernameValidate.Success = false
		isValidRequest = false
	}
	if request.UserInfo.Email == nil || *(request.UserInfo.Email) == "" {
		response.EmailValidate.Success = false
		isValidRequest = false
	}
	if request.Password == nil || *(request.Password) == "" {
		response.PasswordValidate.Success = false
		isValidRequest = false
	}

	if !validation_tools.ValidateEmail(*request.UserInfo.Email) {
		response.EmailValidate.Success = false
		isValidRequest = false
	}

	if !isValidRequest {
		w.WriteHeader(200)
		encoder := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err = encoder.Encode(response)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}


	user := models.User{Username: *request.UserInfo.Username, Email: *request.UserInfo.Email, Password: HashPassword(*request.Password)}
	err = s.St.AddUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	sessionId, err := s.CreateSession(user.ID)
	if err != nil {
		log.Println(err.Error())
	}
	expiration := time.Now().Add(2 * 24 * time.Hour)
	cookie := http.Cookie{Name: "SessionID", Value: sessionId, Expires: expiration}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleLoguot(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	c := http.Cookie{
		Name:   "SessionID",
		MaxAge: 0,
		Value:  "",
	}
	http.SetCookie(w, &c)
	sessionId := cookie.Value
	s.DeleteSession(sessionId)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	var userPerPage int32 = 10
	page, err := strconv.Atoi(r.FormValue("page"))
	users, err := s.St.GetSortedUsers(int32(page)*userPerPage, userPerPage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	usersJson := make([]models.AllInfoJSON, 0)
	for idx, _ := range users {
		userInfo := models.UserInfoJSON{
			Username: &users[idx].Username,
			Email: &users[idx].Email,
			Userpic: &users[idx].Pic,
		}
		if *(userInfo.Username) == "" {
			userInfo.Username = nil
		}
		if *(userInfo.Email) == "" {
			userInfo.Email = nil
		}
		if *(userInfo.Userpic) == "" {
			userInfo.Userpic = nil
		}
		score := strconv.Itoa(int(users[idx].Score))
		gameInfo := models.GameInfoJSON{Score: &score}
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
		&models.ValidateJSON{
			Success: true,
		},
		&models.ErrorJSON{},
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	userId := r.Context().Value("userID").(int32)
	user, err := s.St.GetUserById(userId)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	if request.OldPass != nil {
		if HashPassword(*request.OldPass) != user.Password {
			response.OldPassValidate.Success = false
		} else {
			if request.NewPass != nil {
				err := s.St.ChangePassword(userId, HashPassword(*request.NewPass))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
	}
	if request.UserInfo.Username != nil {
		err := s.St.ChangeUsername(userId, *request.UserInfo.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if request.UserInfo.Userpic != nil {
		err := s.St.ChangePic(userId, *request.UserInfo.Userpic)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if request.UserInfo.Email != nil {
		if !validation_tools.ValidateEmail(*request.UserInfo.Email) {
			response.EmailValidate.Success = false
		} else {
			err := s.St.ChangeEmail(userId, *request.UserInfo.Username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
	if username == "" {
		w.WriteHeader(404)
		return
	}
	user, err := s.St.GetUserByName(username)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	score := strconv.Itoa(int(user.Score))
	userJson := models.AllInfoJSON{
		&models.UserInfoJSON{
			Username: &user.Username,
			Userpic:  &user.Pic,
			Email:    &user.Email,
		},
		&models.GameInfoJSON{
			Score: &score,
		},
	}
	if userJson.UserInfo.Username == nil || *(userJson.UserInfo.Username) == "" {
		userJson.UserInfo.Username = nil
	}
	if userJson.UserInfo.Email == nil || *(userJson.UserInfo.Email) == "" {
		userJson.UserInfo.Email = nil
	}
	if userJson.UserInfo.Userpic == nil || *(userJson.UserInfo.Userpic) == "" {
		userJson.UserInfo.Userpic = nil
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(userJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
