package controllers

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
)

func CreateAccount(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var body models.User
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.MailAlreadyUsed(db, body.Mail) {
		w.WriteHeader(http.StatusOK)
		_, mailErr := w.Write([]byte("Email already in use"))
		if mailErr != nil {
			log.Printf("Failed to write response: %v", mailErr)
		}
		return
	}

	creatingAccountErr := user.CreateAccount(db, body.Username, body.Password, body.Mail, body.Location)
	if creatingAccountErr != nil {
		http.Error(w, "Invalid parameters to create an account", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

var CreateAccountEndPoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "user",
	HandlerFunction: CreateAccount,
}
