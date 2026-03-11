package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type messageResponse struct {
	Messages []interface{} `json:"messages"`
	Count    int           `json:"count"`
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50

	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	response := messageResponse{
		Messages: []interface{}{},
		Count:    0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
