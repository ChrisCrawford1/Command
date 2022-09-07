package handlers

import (
	"encoding/json"
	"github.com/ChrisCrawford1/Command/internal/auth"
	"github.com/ChrisCrawford1/Command/internal/models"
	"github.com/ChrisCrawford1/Command/internal/responses"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type RequestHandler struct {
	Commands interface {
		CreateCommand(creationRequest models.CommandCreationRequest) (bool, error)
		GetByUUID(uuid string) (models.Command, error)
	}
	Users interface {
		GetByEmail(email string) (models.User, error)
		GetByUUID(uuid string) (models.User, error)
	}
}

func (handler *RequestHandler) CreateCommand(w http.ResponseWriter, r *http.Request) {
	var creationRequest models.CommandCreationRequest
	err := json.NewDecoder(r.Body).Decode(&creationRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.Error{Message: err.Error()})
		return
	}

	_, err = handler.Commands.CreateCommand(creationRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.Success{Message: "Command created successfully"})
}

func (handler *RequestHandler) GetCommand(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	foundCommand, _ := handler.Commands.GetByUUID(uuid)

	if foundCommand == (models.Command{}) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responses.Error{Message: "No command could be found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foundCommand.AsDTO())
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
