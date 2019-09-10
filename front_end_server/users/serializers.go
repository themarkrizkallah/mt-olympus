package users

// UserPayload Represents the payload for POST /signup
type UserPayload struct {
	UserName  string `json:"user_name" bson:"user_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Address   string `json:"address" bson:"address"`
}

// LoginPayload Represents the paylod for POST /login
type LoginPayload struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Parse Parses a UserPayload into a User
func (u *UserPayload) Parse() User {
	return User{
		UserName:  u.UserName,
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Address:   u.Address,
	}
}

// Parse Parses a LoginPayload into a User
func (u *LoginPayload) Parse() User {
	return User{
		UserName: u.UserName,
		Email:    u.Email,
		Password: u.Password,
	}
}
