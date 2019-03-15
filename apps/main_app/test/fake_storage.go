package main

import (
	"2019_1_Auteam/models"
	"2019_1_Auteam/server"
	"fmt"
)

type FakeStorage struct {
}

func (st *FakeStorage) AddUser(user *models.User) error {
	return nil
}

func (st *FakeStorage) GetUserByName(username string) (models.User, error) {
	if username == "olzudina" {
		return models.User{Username: "olzudina", Password: server.HashPassword("password")}, nil
	} else {
		return models.User{}, fmt.Errorf("There is no %s", username)
	}
}

func (st *FakeStorage) GetUserById(id int32) (models.User, error) {
	if id < 0 {
		return models.User{}, fmt.Errorf("id < 0")
	}
	return models.User{}, nil
}

func (st *FakeStorage) GetAllUsers() (models.Users, error) {
	return models.Users{}, nil
}

func (st *FakeStorage) GetSortedUsers(from int32, count int32) (models.Users, error) {
	return models.Users{
		models.User{Username: "olzudina", Email: "olzudina@mail.ru", Score: 123456},
		models.User{Username: "ekislukha", Email: "", Score: 12345},
		models.User{Username: "mlozhechko", Email: "mlozhechko@mail.ru", Score: 0},
		models.User{Username: "", Email: "dpoponkin@mail.ru", Score: 123},
		models.User{Username: "vsokolov", Email: "vsokolov@mail.ru", Score: 12},
	}, nil
}

func (st *FakeStorage) ChangeUsername(userID int32, newUsername string) error {
	if userID < 0 {
		return fmt.Errorf("id < 0")
	}
	return nil
}

func (st *FakeStorage) ChangePassword(userID int32, newPassword string) error {
	if userID < 0 {
		return fmt.Errorf("id < 0")
	}
	return nil
}

func (st *FakeStorage) ChangeEmail(userID int32, newEmail string) error {
	if userID < 0 {
		return fmt.Errorf("id < 0")
	}
	return nil
}

func (st *FakeStorage) ChangePic(userID int32, newPic string) error {
	if userID < 0 {
		return fmt.Errorf("id < 0")
	}
	return nil
}

func (st *FakeStorage) UpdateScore(userID int32, newScore int32) error {
	if userID < 0 {
		return fmt.Errorf("id < 0")
	}
	return nil
}

func (st *FakeStorage) UpdateLevel(userID int32, newLevel int32) error {
	if userID < 0 {
		return fmt.Errorf("id < 0")
	}
	return nil
}
