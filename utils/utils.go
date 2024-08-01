package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)


var Validator = validator.New()

type ApiError struct {
	Err string
	Status int
}

type ApiFunc func (http.ResponseWriter, *http.Request) error

func (e ApiError) Error() string {
	return e.Err
}


func ParseJSON(r *http.Request, payload any) error {	
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	WriteJson(w, statusCode, map[string]string{"error": err.Error()})
}

func MakeHTTPHandler(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if e, ok := err.(ApiError); ok {
				WriteJson(w, e.Status, e)
				return
			}
			WriteJson(w, http.StatusInternalServerError, ApiError{Err: "internal server error", Status: http.StatusInternalServerError})
		}
	}
}