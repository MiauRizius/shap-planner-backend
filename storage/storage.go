package storage

import (
	"database/sql"
	"errors"
	"shap-planner-backend/models"

	_ "github.com/glebarez/go-sqlite"
)

var ErrNotFound = sql.ErrNoRows
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
    	username TEXT UNIQUE NOT NULL,
    	password TEXT NOT NULL,
    	role TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}

	//Create refresh token-table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS refresh_tokens(
    	id TEXT PRIMARY KEY,
    	user_id TEXT NOT NULL,
        token_hash TEXT NOT NULL,
        expires_at INTEGER NOT NULL,
        created_at INTEGER NOT NULL,
        revoked INTEGER NOT NULL DEFAULT 0,
        device_info TEXT,
        FOREIGN KEY(user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_refresh_token_hash ON refresh_tokens(token_hash)`)
	if err != nil {
		return err
	}

	//Create Expenses-Table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS expenses(
    	id TEXT PRIMARY KEY
        
	)`)
	return err
}

// Users
func AddUser(user models.User) error {
	_, err := DB.Exec("INSERT INTO users(id, username, password, role) VALUES (?, ?, ?, ?)", user.ID, user.Username, user.Password, user.Role)
	return err
}
func GetUserByUsername(username string) (models.User, error) {
	row := DB.QueryRow("SELECT * FROM users WHERE username = ?", username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	return user, err
}
func GetUserById(id string) (models.User, error) {
	row := DB.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

// Refresh Tokens
func AddRefreshToken(token models.RefreshToken) error {
	_, err := DB.Exec("INSERT INTO refresh_tokens(id, user_id, token_hash, expires_at, created_at, revoked, device_info) VALUES (?, ?, ?, ?, ?, ?, ?)",
		token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt, token.Revoked, token.DeviceInfo)
	return err
}
func GetRefreshToken(token string) (models.RefreshToken, error) {
	row := DB.QueryRow("SELECT * FROM refresh_tokens WHERE token_hash = ?", token)
	var refresh_token models.RefreshToken
	err := row.Scan(&refresh_token.ID, &refresh_token.UserID, &refresh_token.Token, &refresh_token.ExpiresAt, &refresh_token.CreatedAt, &refresh_token.Revoked, &refresh_token.DeviceInfo)
	return refresh_token, err
}
func RevokeRefreshToken(tokenID string) error {
	if DB == nil {
		return errors.New("db not initialized")
	}

	res, err := DB.Exec(`
		UPDATE refresh_tokens
		SET revoked = 1
		WHERE id = ?
	`, tokenID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
func RevokeAllRefreshTokensForUser(userID string) error {
	if DB == nil {
		return errors.New("db not initialized")
	}

	_, err := DB.Exec(`
		UPDATE refresh_tokens
		SET revoked = 1
		WHERE user_id = ?
	`, userID)
	return err
}
