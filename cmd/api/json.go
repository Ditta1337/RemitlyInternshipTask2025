package main

import (
	"encoding/json"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/dto/responses"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func (app *application) writeJSONError(w http.ResponseWriter, status int, message string) {

	if err := writeJSON(w, status, &responses.Error{Error: message}); err != nil {
		app.logger.Errorf("error writing JSOn error: %s", err.Error())
	}
}

func (app *application) writeJSONResponse(w http.ResponseWriter, status int, data any) error {
	return writeJSON(w, status, data)
}
