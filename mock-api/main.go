package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Errors string `json:"errors"`
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/auth/change-password", ChangePassword).Methods(http.MethodPost, http.MethodOptions)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func validate(existingPassword, password, confirmPassword string) (string, bool) {
	if existingPassword == "" || password == "" || confirmPassword == "" {
		return "Missing required field", false
	} else if existingPassword != "Password1" {
		return "Password supplied was incorrect or user is not active", false
	} else if password != confirmPassword {
		return "Confirmation did not match new password", false
	}

	return "", true
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")

	if r.Method == http.MethodOptions {
		return
	}

	existingPassword := r.FormValue("existingPassword")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	errorMessage, ok := validate(existingPassword, password, confirmPassword)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(Response{Errors: errorMessage})
}
