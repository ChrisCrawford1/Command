package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ChrisCrawford1/Command/internal/middleware"
	"github.com/ChrisCrawford1/Command/internal/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockUserModel struct{}

var loggedInUUID = uuid.NewV4()

func (m *MockUserModel) GetByEmail(email string) (models.User, error) {
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

func (m *MockUserModel) GetByUUID(uuid string) (models.User, error) {
	users := map[string]models.User{}

	users[loggedInUUID.String()] = models.User{
		ID:        1,
		UUID:      loggedInUUID,
		Name:      "John",
		Email:     "john@test.com",
		Password:  "$2a$12$.PLe8D00F8qEfHWQVzq8u.7qi397Cy22KaDD5F1Ken97/pgQjk8Qu",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if val, ok := users[uuid]; ok {
		return val, nil
	}

	return models.User{}, nil
}

func TestServer_GetUser(t *testing.T) {
	t.Run("Will get a JWT for stored user with correct credentials with 200", func(t *testing.T) {
		postBody := map[string]interface{}{
			"email":    "john@test.com",
			"password": "password",
		}

		body, _ := json.Marshal(postBody)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(body))

		server := RequestHandler{Users: &MockUserModel{}}

		http.HandlerFunc(server.Login).ServeHTTP(rec, req)

		received := strings.Split(rec.Body.String(), ".")

		if len(received) != 3 {
			t.Errorf("Expected a split token of length 3, received %d", len(received))
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected a response code of 200, received %d", rec.Code)
		}
	})

	t.Run("Will get an error if a stored users password is incorrect", func(t *testing.T) {
		postBody := map[string]interface{}{
			"email":    "john@test.com",
			"password": "wrong-password",
		}

		body, _ := json.Marshal(postBody)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(body))

		server := RequestHandler{Users: &MockUserModel{}}

		http.HandlerFunc(server.Login).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected a response code of 200, received %d", rec.Code)
		}
	})

	t.Run("Will get a not found if user record does not exist", func(t *testing.T) {
		postBody := map[string]interface{}{
			"email":    "dave@test.com",
			"password": "password",
		}

		body, _ := json.Marshal(postBody)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(body))

		server := RequestHandler{Users: &MockUserModel{}}

		http.HandlerFunc(server.Login).ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected a response code of 404, received %d", rec.Code)
		}
	})
}

func TestRequestHandler_GetMe(t *testing.T) {
	t.Run("Will a 200 when fetching the logged in user", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/me", nil)
		if err != nil {
			t.Fatal(err)
		}

		t.Setenv("JWT_SIGN", "INSECURE_SIGN_STRING")
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI1NzM1NTYsInVzZXJJZCI6ImZmMjc3N2QyLWE2NjgtNGIzYS05MDEyLTU0ZmM5NmJjMmNmMiJ9.w5CkYlZ0z4PvBVDoMurL1mijE-9CHJsGeo4OESQcdVA"
		req.Header.Set("Authorization", "Bearer "+validToken)

		ctx := req.Context()
		ctx = context.WithValue(ctx, "userId", loggedInUUID.String())
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		server := RequestHandler{Users: &MockUserModel{}}

		server.GetMe(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected a response code of 200, received %d", rec.Code)
		}
	})

	t.Run("Will a 404 if userId in context doesnt match", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/me", nil)
		if err != nil {
			t.Fatal(err)
		}

		t.Setenv("JWT_SIGN", "INSECURE_SIGN_STRING")
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI1NzM1NTYsInVzZXJJZCI6ImZmMjc3N2QyLWE2NjgtNGIzYS05MDEyLTU0ZmM5NmJjMmNmMiJ9.w5CkYlZ0z4PvBVDoMurL1mijE-9CHJsGeo4OESQcdVA"
		req.Header.Set("Authorization", "Bearer "+validToken)

		ctx := req.Context()
		ctx = context.WithValue(ctx, "userId", "failure-waiting-to-happen")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		server := RequestHandler{Users: &MockUserModel{}}

		server.GetMe(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected a response code of 404, received %d", rec.Code)
		}
	})
}

func TestJwtValidationOnRequest(t *testing.T) {
	t.Run("Will return 401 is authorization header is not present", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/me", nil)
		if err != nil {
			t.Fatal(err)
		}

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		rec := httptest.NewRecorder()
		handler := middleware.ValidateJwtToken(testHandler)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected a response code of 401, received %d", rec.Code)
		}
	})

	t.Run("Will return 401 for an invalid jwt signature", func(t *testing.T) {
		invalidSignatureToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjIxNDU2MjksInVzZXJJZCI6ImZhMjEyYThjLTZjZDYtNDk5MS05MTUwLWQxOWZjZTNhMjRlOSJ9.UWbWsYFnyVqvQAEsR-Lu8Q3QJ29MQsbC_OeDbZumIyw"

		req, err := http.NewRequest("GET", "/users/me", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+invalidSignatureToken)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		rec := httptest.NewRecorder()
		handler := middleware.ValidateJwtToken(testHandler)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected a response code of 401, received %d", rec.Code)
		}
	})

	t.Run("Will return 401 for an expired jwt token", func(t *testing.T) {
		t.Setenv("JWT_SIGN", "INSECURE_SIGN_STRING")
		expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjIyOTkwNTUsInVzZXJJZCI6ImZhMjEyYThjLTZjZDYtNDk5MS05MTUwLWQxOWZjZTNhMjRlOSJ9.OHAnhNKyC2etsq67nEaqlnHBQDBgW1ZsSxiljz4ZAb8"

		req, err := http.NewRequest("GET", "/users/me", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+expiredToken)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})

		rec := httptest.NewRecorder()
		handler := middleware.ValidateJwtToken(testHandler)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected a response code of 401, received %d", rec.Code)
		}
	})
}
