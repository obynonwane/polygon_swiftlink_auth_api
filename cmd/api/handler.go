package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/obynonwane/polygon_swiftlink_auth_api/data"
)

	func (app *Config) Signup(w http.ResponseWriter, r *http.Request) {

		log.Println("reach the auth service")

		var requestPayload data.SignupPayload

		//extract the requestbody
		err := app.readJSON(w, r, &requestPayload)
		if err != nil {
			log.Println(err)
			app.errorJSON(w, err, nil)
			return
		}

		user, err := app.Repo.Signup(requestPayload)
		if err != nil {
			if err == sql.ErrNoRows {
				app.errorJSON(w, errors.New("no record found"), nil, http.StatusBadRequest)
				return
			}

			app.errorJSON(w, err, nil, http.StatusInternalServerError)
			return
		}

		log.Println(user)

		payload := jsonResponse{
			Error:      false,
			StatusCode: http.StatusAccepted,
			Message:    "users retrieved successfully",
			Data:       user,
		}

		app.writeJSON(w, http.StatusAccepted, payload)
	}
func (app *Config) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := app.Repo.GetAll()
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorJSON(w, errors.New("no record found"), nil, http.StatusBadRequest)
			return
		}

		app.errorJSON(w, err, nil, http.StatusInternalServerError)
		return
	}

	log.Println(users)

	payload := jsonResponse{
		Error:      false,
		StatusCode: http.StatusAccepted,
		Message:    "users retrieved successfully",
		Data:       users,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
