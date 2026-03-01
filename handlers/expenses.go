package handlers

import (
	"encoding/json"
	"net/http"
	"shap-planner-backend/models"
	"shap-planner-backend/storage"
	"shap-planner-backend/utils"
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
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if claims.UserID != body.Expense.PayerID { // You cannot create an expense in the name of another user
			http.Error(w, "Invalid request", http.StatusUnauthorized)
			return
		}
		if body.Expense.ID != "" {
			http.Error(w, "Invalid request", http.StatusUnauthorized)
			return
		}
		body.Expense.ID = utils.GenerateUUID()
		for _, share := range body.Shares {
			if share.ID != "" {
				http.Error(w, "Invalid request", http.StatusUnauthorized)
				return
			}
			if share.ExpenseID != "" {
				http.Error(w, "Invalid request", http.StatusUnauthorized)
				return
			}
			share.ExpenseID = body.Expense.ID
			share.ID = utils.GenerateUUID()
			err := storage.AddShare(&share)
			if err != nil {
				http.Error(w, "Error adding expense", http.StatusBadRequest) // Should never happen
				return
			}
		}
		err := storage.AddExpense(&body.Expense)
		if err != nil {
			http.Error(w, "Error adding expense", http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"expense": body.Expense,
			"shares":  body.Shares,
		})
		if err != nil {
			println(err.Error())
			return
		}
		break
	case http.MethodPut: // -> Update Expense
		break
	case http.MethodDelete: // -> Delete Expense
	}
}
func AdminPanel(w http.ResponseWriter, r *http.Request) {}
