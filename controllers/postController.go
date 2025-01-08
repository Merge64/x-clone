package controllers

import (
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
)

func CreateAccount(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var locationAux *string
	username := r.FormValue("username")
	password := r.FormValue("password")
	mail := r.FormValue("mail")
	location := r.FormValue("location")

	if location != constants.EMPTY {
		locationAux = &location
	}

	if username == constants.EMPTY || password == constants.EMPTY || mail == constants.EMPTY {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if user.MailAlreadyUsed(db, mail) {
		w.WriteHeader(http.StatusOK)
		_, mailErr := w.Write([]byte("Email already in use"))
		if mailErr != nil {
			log.Printf("Failed to write response: %v", mailErr)
		}
		return
	}

	creatingAccountErr := user.CreateAccount(db, username, password, mail, locationAux)
	if creatingAccountErr != nil {
		http.Error(w, "Invalid parameters to create an account", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// TODO: In the future implement JWT.
func UserLogin(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	inputUser := r.FormValue("username-or-email")
	password := r.FormValue("password")

	if inputUser == constants.EMPTY || password == constants.EMPTY {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if !user.ValidateCredentials(db, inputUser, password) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("username or password is incorrect"))
		if err != nil {
			log.Printf("Failed to write response: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_, err := w.Write([]byte("Login successful"))
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

var CreateAccountEndPoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "user",
	HandlerFunction: CreateAccount,
}

var UserLoginEndPoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "login",
	HandlerFunction: UserLogin,
}
