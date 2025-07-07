package dto

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSignup struct {
	UserLogin
	Email  string `json:"email"`
	Number string `json:"number"`
}
