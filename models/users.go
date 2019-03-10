package models

type User struct {
	ID int32
	Username string
	Email string
	Password string
	Pic string
	Level int32
	Score int32
}

type Users []User

func (u Users) Len() int {
	return len(u)
}

func (u Users) Less(i, j int) bool {
	return u[i].Score > u[j].Score
}

func (u Users) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
