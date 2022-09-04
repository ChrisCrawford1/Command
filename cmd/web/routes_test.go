package main

import (
	"fmt"
	"github.com/ChrisCrawford1/Command/internal/handlers"
	"github.com/ChrisCrawford1/Command/internal/models"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
)

type mockUserModel struct{}

var loggedInUUID = uuid.NewV4()

func (m *mockUserModel) GetByEmail(email string) (models.User, error) {
	users := map[string]models.User{}

	users["john@test.com"] = models.User{
		ID:        1,
		UUID:      uuid.NewV4(),
		Name:      "John",
		Email:     "john@test.com",
		Password:  "$2a$12$.PLe8D00F8qEfHWQVzq8u.7qi397Cy22KaDD5F1Ken97/pgQjk8Qu",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if val, ok := users[email]; ok {
		return val, nil
	}

	return models.User{}, nil
}

func (m *mockUserModel) GetByUUID(uuid string) (models.User, error) {
	return models.User{
		ID:        1,
		UUID:      loggedInUUID,
		Name:      "John",
		Email:     "john@test.com",
		Password:  "$2a$12$.PLe8D00F8qEfHWQVzq8u.7qi397Cy22KaDD5F1Ken97/pgQjk8Qu",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func TestRoutes(t *testing.T) {
	server := handlers.RequestHandler{Users: &mockUserModel{}}
	mux := Routes(&server)

	switch v := mux.(type) {
	case *chi.Mux:
		// Do nothing
	default:
		t.Error(fmt.Sprintf("%T type is not chi.Mux", v))
	}
}
