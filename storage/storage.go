package storage

import(
    "fmt"
    "2019_1_Auteam/models"
    "database/sql"
    _ "github.com/lib/pq"
)


const initStr = `CREATE TABLE IF NOT EXISTS users(
  id INT PRIMARY KEY,
  username VARCHAR(30)  NOT NULL,
  email VARCHAR(30)  NOT NULL,
  password VARCHAR(120) NOT NULL,
  pic VARCHAR(120) DEFAULT NULL,
  lvl INTEGER DEFAULT 0,
  score INTEGER DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS "users_username_uindex" ON users (username);

CREATE UNIQUE INDEX IF NOT EXISTS "users_score_uindex" ON users (score);
`

const addUserStr = `
INSERT INTO "users" ("username", "email", "password", "pic", "lvl", "score")
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
    db *sql.DB
}

func OpenPostgreStorage(connStr string) (*PostgreStorage, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    _, err = db.Exec(initStr)
    if err != nil {
        return nil, err
    }
    return &PostgreStorage{
        db,
    }, nil
}

func (st *PostgreStorage) AddUser(user* models.User) (error) {
    err := st.db.QueryRow(addUserStr, user.Username, user.Email, user.Password, user.Pic, user.Level, user.Score).Scan(&user.ID)
    return err;
}

func (st *PostgreStorage) GetUserById(userId int32) (models.User, error) {
    var user models.User
    err := st.db.QueryRow(`SELECT * FROM users WHERE id = $1`, userId).Scan(&user)
    fmt.Println(err)
    return user, err
}

func (st *PostgreStorage) GetUserByName(username string) (models.User, error) {
    var user models.User
    err := st.db.QueryRow(`SELECT * FROM users WHERE users.username = $1`, username).Scan(&user)
    fmt.Println(err)
    return user, err
}

func (st *PostgreStorage) GetAllUsers() (models.Users, error) {
    var users models.Users
    res, err := st.db.Query(`SELECT * FROM users`)
    fmt.Println(err)
    if err != nil {
        return users, err
    }
    res.Scan(&users)
    return users, err
}

func (st *PostgreStorage) GetSortedUsers(from int32, count int32) (models.Users, error) {
    var users models.Users
    res, err := st.db.Query(`SELECT * FROM users ORDER BY score DESC LIMIT $2 OFFSET $1`, from, count)
    fmt.Println(err)
    if err != nil {
        return users, err
    }
    res.Scan(&users)
    return users, err
}

func (st *PostgreStorage) ChangeUsername(userID int32, newUsername string) (error) {
    _, err := st.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, newUsername, userID)
    fmt.Println(err)
    return err
}

func (st *PostgreStorage) ChangePassword(userID int32, newPassword string) (error) {
    _, err := st.db.Exec(`UPDATE users SET password = $1 WHERE id = $2`, newPassword, userID)
    fmt.Println(err)
    return err
}

func (st *PostgreStorage) ChangeEmail(userID int32, newEmail string) (error) {
    _, err := st.db.Exec(`UPDATE users SET email = $1 WHERE id = $2`, newEmail, userID)
    fmt.Println(err)
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