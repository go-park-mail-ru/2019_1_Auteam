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


func TestList(t *testing.T) {
	recorder := httptest.NewRecorder()
	userCases := [][4]string {
		{
			"olzudina",
			"olzudina@mail.ru",
			"123456",
		},
		{
			"ekislukha",
			"ekislukha@mail.ru",
			"12345",
		},
		{
			"mlozhechko",
			"mlozhechko@mail.ru",
			"1234",
		},
		{
			"dpoponkin",
			"dpoponkin@mail.ru",
			"123",
		},
		{
			"vsokolov",
			"vsokolov@mail.ru",
			"12",
		},
	}

	expectedUsers := make([]models.AllInfoJSON, 0, 5)
	for idx, _ := range userCases {
		expectedUsers = append(expectedUsers,
			models.AllInfoJSON{
				&models.UserInfoJSON {
					Username: &(userCases[idx][0]),
					Email:  &(userCases[idx][1]),
				},
				&models.GameInfoJSON {
					Score:  &(userCases[idx][2]),
				},
			})
	}

	req, _ := http.NewRequest("GET", "/user/list", nil)
	router.ServeHTTP(recorder, req)
	if recorder.Code != 200 {
		t.Errorf("Expected response with status 200, actual - %v", recorder.Code)
	}

	expected, _ := json.Marshal(expectedUsers)
	assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body differs")
}

func TestUsernameGood(t* testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/username", nil)
	router.ServeHTTP(recorder, req)
}