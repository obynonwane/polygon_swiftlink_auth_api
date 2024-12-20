package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/obynonwane/polygon_swiftlink_auth_api/data"
	"github.com/obynonwane/polygon_swiftlink_auth_api/util"
)

type loginUserResponse struct {
	AccessToken string             `json:"access_token"`
	User        CreateUserResponse `json:"user"`
}

type CreateUserResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (app *Config) VerifyToken(w *http.ResponseWriter, r http.Request) {
	// make a call to retrive the user
}

func (app *Config) Signup(w http.ResponseWriter, r *http.Request) {

	var requestPayload data.SignupPayload

	//extract the requestbody
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
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

	payload := jsonResponse{
		Error:      false,
		StatusCode: http.StatusAccepted,
		Message:    "user account created successfully",
		Data:       user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {

	var req data.LoginUserRequest

	err := app.readJSON(w, r, &req)
	if err != nil {
		app.errorJSON(w, errors.New("error while reading inputs"), http.StatusBadRequest)
		return
	}

	//get user and compare password
	user, err := app.Repo.GetUserWithEmail(req.Email)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, errors.New("error signing-in"), http.StatusBadRequest)
		return

	}

	//compare user password
	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		log.Println("Error 2")
		app.errorJSON(w, errors.New("user unauthorised"), http.StatusUnauthorized)
		return
	}

	//create access token //AccessTokenDuration: Needed to be in .env file
	pasetoDetail := &TokenType{
		AccessTokenDuration: 15 * time.Minute,
	}
	accessToken, err := app.tokenMaker.CreateToken(req.Email, pasetoDetail.AccessTokenDuration)
	if err != nil {
		app.errorJSON(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	rsp := CreateUserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
	}

	rep2 := loginUserResponse{
		AccessToken: accessToken,
		User:        rsp,
	}

	payload := jsonResponse{
		Error:   false,
		Message: "login successful",
		Data:    rep2,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
func (app *Config) GetUsers(w http.ResponseWriter, r *http.Request) {

	log.Println("welcome here it reached here")

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
