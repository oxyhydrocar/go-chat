package handlers

import (
	"encoding/json"
	"net/http"

	"go-chat/models"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	hub := models.GetHub()
	users := hub.GetActiveUsers()

	response := make([]map[string]string, 0, len(users))
	for _, user := range users {
		response = append(response, map[string]string{
			"id":       user.ID,
			"username": user.Username,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
