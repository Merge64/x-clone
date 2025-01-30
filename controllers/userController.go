package controllers

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if !user.IsEmail(mail) {
		w.WriteHeader(http.StatusOK)
		_, mailErr := w.Write([]byte("Invalid email"))
		if mailErr != nil {
			log.Printf("Failed to write response: %v", mailErr)
		}
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

	if user.UsernameAlreadyUsed(db, username) {
		w.WriteHeader(http.StatusOK)
		_, usernameErr := w.Write([]byte("Username already in use"))
		if usernameErr != nil {
			log.Printf("Failed to write response: %v", usernameErr)
		}
		return
	}

	creatingAccountErr := user.CreateAccount(db, username, password, mail, locationAux)
	if creatingAccountErr != nil {
		http.Error(w, "Invalid parameters to create an account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("Account created successfully"))
	if err != nil {
		return
	}
}

// TODO: In the future implement JWT.
func UserLoginHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

func FollowUserHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	followingID, getIDErr := getUserID(r)
	if getIDErr != nil {
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}

	followedUserID, atoiErr := strconv.Atoi(r.PathValue("userid"))
	if atoiErr != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if followErr := user.FollowAccount(db, followingID, uint(followedUserID)); followErr != nil {
		fmt.Println(followErr)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Follows user successfully"))
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func UnfollowUserHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	followingID, getIDErr := getUserID(r)
	if getIDErr != nil {
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}

	followedUserID, atoiErr := strconv.Atoi(r.PathValue("userid"))
	if atoiErr != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if unfollowErr := user.UnfollowAccount(db, followingID, uint(followedUserID)); unfollowErr != nil {
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Unfollows user successfully"))
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func getUserID(r *http.Request) (uint, error) {
	var currentUser models.User
	if decodeErr := json.NewDecoder(r.Body).Decode(&currentUser); decodeErr != nil {
		return 0, decodeErr
	}
	return currentUser.ID, nil
}

var FollowUserEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "follow/{userid}",
	HandlerFunction: FollowUserHandler,
}

var UnfollowUserEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.BASEURL + "unfollow/{userid}",
	HandlerFunction: UnfollowUserHandler,
}

var UserSignUpEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "signup",
	HandlerFunction: SignUpHandler,
}

var UserLoginEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "login",
	HandlerFunction: UserLoginHandler,
}
