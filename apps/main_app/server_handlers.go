package main

import (
	"2019_1_Auteam/models"
	"2019_1_Auteam/validation_tools"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

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
	if request.UserInfo.Username != nil ||
		request.UserInfo.Email != nil ||
		request.Password != nil ||
		!validation_tools.ValidateEmail(*request.UserInfo.Email) {
		w.WriteHeader(400)
		return
	}
	user := models.User{Username: *request.UserInfo.Username, Email: *request.UserInfo.Email, Password: hashPassword(*request.Password)}
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
	users, err := s.st.GetSortedUsers(int32(page)*userPerPage, userPerPage)
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
		if !validation_tools.ValidateEmail(*request.UserInfo.Email) {
			response.EmailValidate.Success = false
		} else {
			err := s.st.ChangeEmail(userId, *request.UserInfo.Username)
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
	user, err := s.st.GetUserByName(username)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	userJson := models.AllInfoJSON{
		&models.UserInfoJSON{
			Username: &user.Username,
			Userpic:  &user.Pic,
			Email:    &user.Email,
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
