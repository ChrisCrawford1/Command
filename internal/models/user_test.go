package models

import (
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
)

func TestUser_AsDTO(t *testing.T) {
	t.Run("Will have expected fields set not empty", func(t *testing.T) {

		user := User{
			ID:        1,
			UUID:      uuid.NewV4(),
			Name:      "John",
			Email:     "john@test.com",
			Password:  "supersecurepassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		setAsDto := user.AsDTO()

		if setAsDto.UUID.String() == "" {
			t.Error("User DTO uuid was empty, value was expected")
		}

		if setAsDto.Name == "" {
			t.Error("User DTO name was empty, value was expected")
		}

		if setAsDto.Email == "" {
			t.Error("User DTO email was empty, value was expected")
		}
	})
}
