package model

type Signup struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Login struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
