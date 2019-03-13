package models

type User struct {
	ID int32 `db:"id"`
	Username string `db:"username"`
	Email string `db:"email"`
	Password string `db:"password"`
	Pic string `db:"pic"`
	Level int32 `db:"lvl"`
	Score int32 `db:"score"`
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
