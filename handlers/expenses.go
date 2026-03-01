package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shap-planner-backend/models"
	"shap-planner-backend/storage"
	"shap-planner-backend/utils"
	"time"
)

func Expenses(w http.ResponseWriter, r *http.Request) {
	claims, _ := utils.IsLoggedIn(w, r)

	switch r.Method {
	case http.MethodGet: // -> Get Expenses
		break
	case http.MethodPost: // -> Create Expense
		var body struct {
			Expense models.Expense        `json:"expense"`
			Shares  []models.ExpenseShare `json:"shares"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Println("POST [api/expense] " + r.RemoteAddr + ": " + err.Error())
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if claims.UserID != body.Expense.PayerID { // You cannot create an expense in the name of another user
			log.Println("POST [api/expense] " + r.RemoteAddr + ": claims.UserID and expense.PayerID does not match")
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		// Set ExpenseID
		if body.Expense.ID != "" {
			log.Println("POST [api/expense] " + r.RemoteAddr + ": Expense ID must be empty")
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		body.Expense.ID = utils.GenerateUUID()
		if body.Expense.CreatedAt == 0 {
			body.Expense.CreatedAt = time.Now().Unix()
		}
		if body.Expense.Amount <= 0 {
			log.Println("POST [api/expense] " + r.RemoteAddr + ": Amount must be greater than zero")
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		// Set ShareIDs and save them
		for _, share := range body.Shares {
			if share.ID != "" {
				log.Println("POST [api/expense] " + r.RemoteAddr + ": Share ID must be empty")
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			if share.ExpenseID != "" {
				log.Println("POST [api/expense] " + r.RemoteAddr + ": Expense ID of Share must be empty")
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			share.ExpenseID = body.Expense.ID
			share.ID = utils.GenerateUUID()
			err := storage.AddShare(&share)
			if err != nil {
				log.Println("POST [api/expense] " + r.RemoteAddr + ": " + err.Error())
				http.Error(w, "Error adding expense", http.StatusBadRequest) // Should never happen
				return
			}
		}
		err := storage.AddExpense(&body.Expense)
		if err != nil {
			log.Println("POST [api/expense] " + r.RemoteAddr + ": " + err.Error())
			http.Error(w, "Error adding expense", http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"expense": body.Expense,
			"shares":  body.Shares,
		})
		if err != nil {
			log.Println("POST [api/expense] " + r.RemoteAddr + ": " + err.Error())
			return
		}
		log.Println("POST [api/expense] " + r.RemoteAddr + ": Successfully added expense and its shares")
		break
	case http.MethodPut: // -> Update Expense
		break
	case http.MethodDelete: // -> Delete Expense
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
func AdminPanel(w http.ResponseWriter, r *http.Request) {}
