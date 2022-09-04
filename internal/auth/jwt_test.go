package auth

import (
	"github.com/ChrisCrawford1/Command/internal/models"
	uuid "github.com/satori/go.uuid"
	"strings"
	"testing"
	"time"
)

func TestGenerateAccessToken(t *testing.T) {
	t.Run("Will return a well formed token", func(t *testing.T) {
		model := models.User{
			ID:        1,
			UUID:      uuid.NewV4(),
			Name:      "Dave",
			Email:     "dave@test.com",
			Password:  "j3h423#34#51234!3&£%",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		generatedToken := GenerateAccessToken(model)
		splitToken := strings.Split(generatedToken, ".")

		if len(splitToken) != 3 {
			t.Errorf("Expected a split token of length 3, received %d", len(splitToken))
		}
	})
}

func TestValidateAccessToken(t *testing.T) {
	t.Run("Will return claims for a well formed token", func(t *testing.T) {
		model := models.User{
			ID:        1,
			UUID:      uuid.NewV4(),
			Name:      "Dave",
			Email:     "dave@test.com",
			Password:  "j3h423#34#51234!3&£%",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		generatedToken := GenerateAccessToken(model)

		valid, claims, err := ValidateAccessToken(generatedToken)

		if !valid {
			t.Error("Expected to receive true for valid token, received false")
		}

		if err != nil {
			t.Error("Expected to receive no error, received error")
		}

		userId, ok := claims["userId"]

		if !ok {
			t.Error("Expected to recive a userId in claims, received none")
		}

		if userId != model.UUID.String() {
			t.Error("Expected user id in claims to match user, they did not")
		}
	})
}
