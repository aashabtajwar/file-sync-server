package users

type RegisterUser struct {
	Id       string
	First    string
	Last     string
	Email    string
	Password string
}

type LoginUser struct {
	Email    string
	Password string
}
