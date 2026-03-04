package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Expense struct {
	ID            string   `json:"id"`
	PayerID       string   `json:"payer_id"`
	Amount        int64    `json:"amount"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Attachments   []string `json:"attachments"`
	CreatedAt     int64    `json:"created_at"`
	LastUpdatedAt int64    `json:"last_updated_at"`
}

type ExpenseShare struct {
	ID         string `json:"id"`
	ExpenseID  string `json:"expense_id"`
	UserID     string `json:"user_id"`
	ShareCents int64  `json:"share_cents"`
}
