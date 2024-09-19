package token

import "time"

// manages the creation and verification of the tokens
type Maker interface {
	//CreateToken creates a new token for specific username and duration
	CreateToken(email string, duration time.Duration) (string, error)

	//VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
