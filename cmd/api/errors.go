package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("internal server error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	app.writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request response: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	app.writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found response: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	app.writeJSONError(w, http.StatusNotFound, "resource not found")
}
