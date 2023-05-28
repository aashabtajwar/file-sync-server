package users

type RegisterUser struct {
	First    string
	Last     string
	Email    string
	Password string
}

type LoginUser struct {
	Email    string
	Password string
}
