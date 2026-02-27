package storage

import (
	"database/sql"
	"shap-planner-backend/models"

	_ "github.com/glebarez/go-sqlite"
)

var DB *sql.DB

func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite", filepath)
	if err != nil {
		return err
	}

	//Create Users-Table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users(
    	id TEXT PRIMARY KEY,
    	username TEXT UNIQUE,
    	password TEXT
	);`)
	if err != nil {
		return err
	}

	//Create Expenses-Table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS expenses(
    	id TEXT PRIMARY KEY
        
	)`)
	return err
}

func AddUser(user models.User) error {
	_, err := DB.Exec("INSERT INTO users(id, username, password, role) VALUES (?, ?, ?, ?)", user.ID, user.Username, user.Password, user.Role)
	return err
}

func GetUserByUsername(username string) (models.User, error) {
	row := DB.QueryRow("SELECT * FROM users WHERE username = ?", username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

func GetUserById(id string) (models.User, error) {
	row := DB.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}
