package controllers

import (
	"encoding/json"
	"main/utils"
	"net/http"
)

func HandlerDicGetWords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, _ := utils.GetUserClaims(r)
	userId := claims["id"].(float64)

	words, err := utils.GetWordsByUserId(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(words); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
