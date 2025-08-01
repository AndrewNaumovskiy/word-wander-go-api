package main

import (
	"net/http"

	"main/controllers"
	"main/utils"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/signup", controllers.HandlerSignup)
	mux.HandleFunc("/auth/login", controllers.HandlerLogin)

	dicGetWords := http.HandlerFunc(controllers.HandlerDicGetWords)
	mux.Handle("/dictionary/get-words", utils.JWTMiddleware(dicGetWords))

	traiGetAmountWords := http.HandlerFunc(controllers.HandlerTrainGetAmountWords)
	mux.Handle("/training/get-amount-words-for-trainings", utils.JWTMiddleware(traiGetAmountWords))

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler(mux)

	http.ListenAndServe(":3001", handler)
}

type DbUserModel struct {
	Id               int
	MongoId          string
	Email            string
	Password         string
	RegistrationDate string
}
