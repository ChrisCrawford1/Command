package handlers

import (
	"encoding/json"
	"github.com/ChrisCrawford1/Command/internal/auth"
	"github.com/ChrisCrawford1/Command/internal/models"
	"github.com/ChrisCrawford1/Command/internal/responses"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type RequestHandler struct {
	Users interface {
		GetByEmail(email string) (models.User, error)
		GetByUUID(uuid string) (models.User, error)
	}
}

func (handler *RequestHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login models.UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&login)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.Error{Message: err.Error()})
		return
	}

	foundUser, err := handler.Users.GetByEmail(login.Email)

	if foundUser == (models.User{}) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responses.Error{Message: "No matching credentials could be found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(login.Password))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(responses.Error{Message: "Invalid credentials"})
		return
	}

	accessToken := auth.GenerateAccessToken(foundUser)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(accessToken))
}

func (handler *RequestHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	foundUser, _ := handler.Users.GetByUUID(r.Context().Value("userId").(string))

	if foundUser == (models.User{}) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responses.Error{Message: "User could not be found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foundUser.AsDTO())
}
