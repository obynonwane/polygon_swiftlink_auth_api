package data

type SignupPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}
type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
