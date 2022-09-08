package responses

import (
	"errors"
	"github.com/go-playground/validator"
)

type Error struct {
	Message string `json:"Message"`
}

type Success struct {
	Message string `json:"Message"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func MessageForTag(tag string) string {
	switch tag {
	case "required":
		return "Field is required"
	}
	return ""
}

func ProcessErrors(err error) ([]ValidationError, error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		responseErrors := make([]ValidationError, len(ve))
		for i, fe := range ve {
			responseErrors[i] = ValidationError{
				Field:   fe.Field(),
				Message: MessageForTag(fe.Tag()),
			}
		}

		return responseErrors, nil
	} else {
		return nil, err
	}
}
