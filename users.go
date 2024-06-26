package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var errEmailRequired = errors.New("email cannot be empty")
var errFirstNameRequired = errors.New("first cannot be empty")
var errLastNameRequired = errors.New("last cannot be empty")
var errPasswordRequired = errors.New("password cannot be empty")

type UserService struct {
	store Store
}

func NewUserService(s Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes (r *mux.Router){
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
}

func (s *UserService) handleUserRegister (w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *User
	err = json.Unmarshal(body, &payload)
	
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := ValidateUserPayLoad(payload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hashedPassword, err := HashPassword(payload.Password)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating user"})
		return
	}

	payload.Password = hashedPassword

	u, err := s.store.CreateUser(payload)
 
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating user"})
		return
	}

	token, err := createAndSetAuthCookie(u.ID, w)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating session"})
		return
	}


	message := "Here is your token: " + token
	WriteJSON(w, http.StatusCreated, message)
	fmt.Println("JWT:", token)
}

func ValidateUserPayLoad(user *User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.FirstName == "" {
		return errFirstNameRequired
	}

	if user.LastName == "" {
		return errLastNameRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(Envs.JWSecret)
	token, err := CreateJWT(secret, id)

	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name: "Authorisation",
		Value: token,
	})

	return token, nil
}