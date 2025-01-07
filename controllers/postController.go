package controllers

import (
	"encoding/json"
	"gorm.io/gorm"
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

	currentAccount, creatingAccountErr := user.CreateAccount(db, body.Username, body.Password, body.Mail, body.Location)
	if creatingAccountErr != nil {
		http.Error(w, "Invalid course ID format", http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(currentAccount)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)

	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

var CreateAccountEndPoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "/profiles/{Username}",
	HandlerFunction: CreateAccount,
}
