package server

import (
	"2019_1_Auteam/models"
	"2019_1_Auteam/validation_tools"
	"encoding/json"
	"net/http"
	"strconv"
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
		return
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

	sessionId, err := s.CreateSession(user.ID)
	if err != nil {
		log.Println(err.Error())
	}
	SetSessionCoockie(sessionId, w)
	w.WriteHeader(http.StatusOK)
}


func (s *Server) handleLoguot(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
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
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userId := r.Context().Value("userID").(int32)
	user, err := s.St.GetUserById(userId)
	if err != nil {
		log.Println(err.Error())
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
					log.Println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
	}
	if request.UserInfo.Username != nil {
		err := s.St.ChangeUsername(userId, *request.UserInfo.Username)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if request.UserInfo.Userpic != nil {
		err := s.St.ChangePic(userId, *request.UserInfo.Userpic)
		if err != nil {
			log.Println(err.Error())
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
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
