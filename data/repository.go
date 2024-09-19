package data

type Repository interface {
	GetAll() ([]*User, error)
	Signup(payload SignupPayload) (*User, error)

	GetUserWithEmail(email string) (*User, error)
}
