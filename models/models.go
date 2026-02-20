package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Expense struct {
	ID          string `json:"id"`
	Amount      int    `json:"amt"`
	Description string `json:"desc"`

	Payer   User   `json:"payer"`
	Debtors []User `json:"debtors"`
}
