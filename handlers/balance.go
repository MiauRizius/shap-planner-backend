package handlers

import "net/http"

func GetBalance(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userParam := query.Get("user")

	if userParam == "" {

	}
}
