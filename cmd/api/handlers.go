package main

import (
	"net/http"
)

func (app *Config) Signup(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit the auth service",
		Data:    nil,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
