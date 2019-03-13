package storage

import(
    "log"
    "fmt"
    "2019_1_Auteam/models"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

const addUserStr = `
INSERT INTO users (username, email, password, pic, lvl, score)
VALUES
($1, $2, $3, $4, $5, $6)
RETURNING id;
`

type StorageI interface {
    AddUser(user* models.User) (error)
    GetUserByName(username string) (models.User, error)
    GetUserById(id int32) (models.User, error)
    GetAllUsers() (models.Users, error)
    GetSortedUsers(from int32, count int32) (models.Users, error)
    ChangeUsername(userID int32, newUsername string) (error)
    ChangePassword(userID int32, newPassword string) (error)
    ChangeEmail(userID int32, newEmail string) (error)
    ChangePic(userID int32, newPic string) (error)
    UpdateScore(userID int32, newScore int32) (error)
    UpdateLevel(userID int32, newLevel int32) (error)
}

type PostgreStorage struct {
    db *sqlx.DB
}

func OpenPostgreStorage(host string, user string, password string, dbname string) (*PostgreStorage, error) {
    db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname))
    if err != nil {
        return nil, err
    }
    err = db.Ping()
    if err != nil {
        return nil, err
    }
//     db.Exec(`DROP TABLE users
// CREATE TABLE IF NOT EXISTS users(
//   id SERIAL PRIMARY KEY,
//   username VARCHAR(30)  NOT NULL,
//   email VARCHAR(30)  NOT NULL,
//   password VARCHAR(120) NOT NULL,
//   pic VARCHAR(120) DEFAULT NULL,
//   lvl INTEGER DEFAULT 0,
//   score INTEGER DEFAULT 0
// );`)
    return &PostgreStorage{db,}, nil
}

func (st *PostgreStorage) AddUser(user* models.User) (error) {
    err := st.db.QueryRow(addUserStr, user.Username, user.Email, user.Password, user.Pic, user.Level, user.Score).Scan(&user.ID)
    return err;
}

func (st *PostgreStorage) GetUserById(userId int32) (models.User, error) {
    var user models.User
    err := st.db.Get(&user, `SELECT * FROM users WHERE id = $1`, userId)
    log.Println(err)
    return user, err
}

func (st *PostgreStorage) GetUserByName(username string) (models.User, error) {
    var user models.User
    err := st.db.Get(&user, `SELECT * FROM users WHERE users.username = $1`, username)
    log.Println(err)
    return user, err
}

func (st *PostgreStorage) GetAllUsers() (models.Users, error) {
    var users models.Users
    err := st.db.Select(&users, `SELECT * FROM users`)
    log.Println(err)
    if err != nil {
        return users, err
    }
    return users, err
}

func (st *PostgreStorage) GetSortedUsers(from int32, count int32) (models.Users, error) {
    var users models.Users
    err := st.db.Select(&users, `SELECT * FROM users ORDER BY score DESC LIMIT $2 OFFSET $1`, from, count)
    log.Println(err)
    if err != nil {
        return users, err
    }
    return users, nil
}

func (st *PostgreStorage) ChangeUsername(userID int32, newUsername string) (error) {
    _, err := st.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, newUsername, userID)
    log.Println(err)
    return err
}

func (st *PostgreStorage) ChangePassword(userID int32, newPassword string) (error) {
    _, err := st.db.Exec(`UPDATE users SET password = $1 WHERE id = $2`, newPassword, userID)
    log.Println(err)
    return err
}

func (st *PostgreStorage) ChangeEmail(userID int32, newEmail string) (error) {
    _, err := st.db.Exec(`UPDATE users SET email = $1 WHERE id = $2`, newEmail, userID)
    log.Println(err)
    return err
}

func (st *PostgreStorage) ChangePic(userID int32, newPic string) (error) {
    _, err := st.db.Exec(`UPDATE users SET pic = $1 WHERE id = $2`, newPic, userID)
    return err
}

func (st *PostgreStorage) UpdateScore(userID int32, newScore int32) (error) {
    _, err := st.db.Exec(`UPDATE users SET score = $1 WHERE id = $2`, newScore, userID)
    return err
}

func (st *PostgreStorage) UpdateLevel(userID int32, newLevel int32) (error) {
    _, err := st.db.Exec(`UPDATE users SET lvl = $1 WHERE id = $2`, newLevel, userID)
    return err
}