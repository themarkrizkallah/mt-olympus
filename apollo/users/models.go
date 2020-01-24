package users

import "time"

const table = "users"

type User struct {
	Id        string
	Email     string
	Password  string
	CreatedAt time.Time
}

// LoginPayload Represents the paylod for POST /login
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupPayload Represents the paylod for POST /signup
type SignupPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
