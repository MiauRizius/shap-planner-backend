package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"shap-planner-backend/models"
	"strings"

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
    	id TEXT PRIMARY KEY,
        payer_id TEXT NOT NULL,
        amount_cents INTEGER NOT NULL,
        title TEXT NOT NULL,
        description TEXT,
        attachments TEXT,
        created_at INTEGER NOT NULL,
        last_updated_at INTEGER NOT NULL,
        FOREIGN KEY(payer_id) REFERENCES users(id)
	)`)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS expense_shares(
    	id TEXT PRIMARY KEY,
		expense_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		share_cents INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_shares_expense ON expense_shares(expense_id)`)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_shares_user ON expense_shares(user_id)`)
	return err
}

// Expenses
func AddExpense(expense *models.Expense) error {
	var attachmentsData interface{}
	if len(expense.Attachments) > 0 {
		jsonBytes, err := json.Marshal(expense.Attachments)
		if err != nil {
			return err
		}
		attachmentsData = string(jsonBytes)
	} else {
		attachmentsData = nil
	}
	_, err := DB.Exec(`INSERT INTO expenses(id, payer_id, amount_cents, title, description, attachments, created_at, last_updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		expense.ID,
		expense.PayerID,
		expense.Amount,
		expense.Title,
		expense.Description,
		attachmentsData,
		expense.CreatedAt,
		expense.LastUpdatedAt)
	return err
}
func UpdateExpense(expense *models.Expense) error {
	return nil
}
func DeleteExpense(expense *models.Expense) error {
	return nil
}

//	func GetExpenseById(id string) (models.Expense, error) {
//		return nil, nil
//	}
func GetExpensesByUserId(userId string) ([]models.Expense, error) {
	return nil, nil
}
func GetAllExpenses() ([]models.Expense, error) {
	return nil, nil
}

// Expense Shares
func AddShare(share *models.ExpenseShare) error {
	_, err := DB.Exec("INSERT INTO expense_shares(id, expense_id, user_id, share_cents) VALUES (?, ?, ?, ?)",
		share.ID,
		share.ExpenseID,
		share.UserID,
		share.ShareCents)
	return err
}

// Users
func AddUser(user *models.User) error {
	_, err := DB.Exec("INSERT INTO users(id, username, password, role) VALUES (?, ?, ?, ?)", user.ID, strings.ToLower(user.Username), user.Password, user.Role)
	return err
}
func GetUserByUsername(username string) (models.User, error) {
	row := DB.QueryRow("SELECT * FROM users WHERE username = ?", strings.ToLower(username))
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	return user, err
}
func GetUserById(id string) (models.User, error) {
	row := DB.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	return user, err
}

// Refresh Tokens
func AddRefreshToken(token *models.RefreshToken) error {
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
