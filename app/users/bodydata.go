package users

type RegisterUser struct {
	Id        string
	First     string
	Last      string
	Email     string
	Username  string
	Password  string
	CreatedAt string
	UpdatedAt string
}

type LoginUser struct {
	Email    string
	Password string
}
