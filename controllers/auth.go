package controllers

import (
	"encoding/json"
	"net/http"

	"main/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserReqModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseModel struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignupResponseModel struct {
	Message SignupResponseMsgModel `json:"message"`
}
type SignupResponseMsgModel struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func HandlerSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user UserReqModel
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	foundUser, err := utils.GetUserByEmail(user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if foundUser != nil {
		resp := SignupResponseModel{
			Message: SignupResponseMsgModel{
				Message: "User already exists",
				Error:   "Bad Request",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	// TOOD: implement later

	w.Write([]byte("Hello, World!"))
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var user UserReqModel
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	foundUser, err := utils.GetUserByEmail(user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	passValid := ComparePasswords(foundUser.Password, user.Password)
	if foundUser == nil || !passValid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(foundUser)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	resp := LoginResponseModel{
		AccessToken:  token,
		RefreshToken: "", // TODO: implement refresh token
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func ComparePasswords(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
