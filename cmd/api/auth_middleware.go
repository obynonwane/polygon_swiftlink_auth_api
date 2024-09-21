package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/obynonwane/polygon_swiftlink_auth_api/token"
)

type contextKey string

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

// Use the custom type for the context key
const authorizationPayloadKey contextKey = "authorization_payload"

func (app *Config) HandleVerifyToken(w http.ResponseWriter, r *http.Request) {

	// This is a custom handler that combines the authMiddleware and Signin logic
	// Apply authMiddleware
	authMiddleware(app.tokenMaker)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the Signin handler here
		app.verifyToken(w, r)
	})).ServeHTTP(w, r)
}

// AuthMiddleware creates a Chi middleware for authorization
func authMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			app := Config{}
			authorizationHeader := r.Header.Get(authorizationHeaderKey)

			if len(authorizationHeader) == 0 {
				err := errors.New("authorization header is not provided")
				app.errorJSON(w, err, nil, http.StatusUnauthorized)
				return
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				err := errors.New("invalid authorization header format")
				app.errorJSON(w, err, nil, http.StatusUnauthorized)
				return
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != authorizationTypeBearer {
				err := fmt.Errorf("unsupported authorization type %s", authorizationType)
				app.errorJSON(w, err, nil, http.StatusUnauthorized)
				return
			}

			accessToken := fields[1]
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				app.errorJSON(w, err, nil, http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, authorizationPayloadKey, payload)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func (app *Config) verifyToken(w http.ResponseWriter, r *http.Request) {

	authPayload, _ := r.Context().Value(authorizationPayloadKey).(*token.Payload)

	user, err := app.Repo.GetUserWithEmail(authPayload.Email)

	if err != nil {
		if err == sql.ErrNoRows {

			app.errorJSON(w, errors.New("error no user detail retrieved"), http.StatusUnauthorized)
			return
		}

		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:      false,
		Message:    "token verified succesfully",
		Data:       user,
		StatusCode: http.StatusAccepted,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// Define a function to format error responses
// func errorResponse(err error) string {
// 	return fmt.Sprintf(`{"error": "%s"}`, err.Error())
// }
