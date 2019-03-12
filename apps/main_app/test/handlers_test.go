package main

import(
	"2019_1_Auteam/models"
	"2019_1_Auteam/server"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/json"
	"bytes"
)

var srv = &server.Server{&FakeStorage{}, nil}
var router = server.CreateRouter(srv)

func TestLoginQueryGood(t *testing.T) {
	recorder := httptest.NewRecorder()
	username := "olzudina"
	password := "password"
	jsmodel := models.LoginRequestJSON {
		Username: &username,
		Password: &password,
	}
	jsonString, _ := json.Marshal(jsmodel)
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewReader(jsonString))
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected response with status OK, actual - %v", recorder.Code)
	}
}

func TestLoginQueryBad(t *testing.T) {
	recorder := httptest.NewRecorder()
	username := "olzudina"
	password := "passwor1"
	jsmodel := models.LoginRequestJSON {
		Username: &username,
		Password: &password,
	}
	jsonString, _ := json.Marshal(jsmodel)
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewReader(jsonString))
	router.ServeHTTP(recorder, req)
	if recorder.Code != 400 {
		t.Errorf("Expected response with status 400, actual - %v", recorder.Code)
	}
}

func TestSignupGood(t *testing.T) {
	recorder := httptest.NewRecorder()
	username := "olzudina"
	email := "olzudina@example.com"
	password := "passwor1"
	jsmodel := models.SignUpRequestJSON {
		UserInfo: &models.UserInfoJSON {
			Username: &username,
			Email: &email,
		},
		Password: &password,
	}
	jsonString, _ := json.Marshal(jsmodel)
	req, _ := http.NewRequest("POST", "/user/signup", bytes.NewReader(jsonString))
	router.ServeHTTP(recorder, req)
	if recorder.Code != 200 {
		t.Errorf("Expected response with status 200, actual - %v", recorder.Code)
	}
}


func TestSignupBad(t *testing.T) {
	recorder := httptest.NewRecorder()
	username := "olzudina"
	email := "olzudinaqwe"
	password := "passwor1"
	jsmodel := models.SignUpRequestJSON {
		UserInfo: &models.UserInfoJSON {
			Username: &username,
			Email: &email,
		},
		Password: &password,
	}
	jsonString, _ := json.Marshal(jsmodel)
	req, _ := http.NewRequest("POST", "/user/signup", bytes.NewReader(jsonString))
	router.ServeHTTP(recorder, req)
	if recorder.Code != 200 {
		t.Errorf("Expected response with status 200, actual - %v", recorder.Code)
	}

	expectedJSON := models.SignUpResponseJSON {
		&models.ValidateJSON {
			Success: true,
		},
		&models.ValidateJSON {
			Success: true,
		},
		&models.ValidateJSON {
			Success: false,
		},
		&models.ValidateJSON {
			Success: true,
		},
		&models.ErrorJSON {},
	}
	expected, _ := json.Marshal(expectedJSON)
	assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body differs")
}
