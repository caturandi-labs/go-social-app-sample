package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	_ = Validate.RegisterValidation("exists_email", emailExists)

}

func emailExists(fl validator.FieldLevel) bool {
	var registeredEmails = []string{"user@example.com", "admin@example.com"}
	email := fl.Field().String()
	for _, e := range registeredEmails {
		if strings.EqualFold(e, email) {
			return true // Email exists
		}
	}
	return false
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Message string `json:"message"`
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(envelope{Message: message})
}

func writeValidationJSONError(w http.ResponseWriter, status int, message string, errors any) error {
	type envelope struct {
		Message string `json:"message"`
		Errors  any    `json:"errors,omitempty"`
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(envelope{Message: message, Errors: errors})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}

func formatValidationErrors(err error) map[string]string {
	errFields := make(map[string]string)
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, err := range validationErrs {
			errField := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errFields[errField] = fmt.Sprintf("%s is required", errField)
			case "email":
				errFields[errField] = "Invalid email format"
			case "min":
				errFields[errField] = fmt.Sprintf("%s must be greater than %v characters", errField, err.Param())
			case "max":
				errFields[errField] = fmt.Sprintf("%s must be less than %v characters", errField, err.Param())
			case "unique_email":
				errFields[errField] = "Email already taken"
			case "exists_email":
				errFields[errField] = "Email does not exist"
			default:
				errFields[errField] = "Invalid value"
			}
		}
	}
	return errFields
}
