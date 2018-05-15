package models

// Token is a type that holds a JWT inside and the user that owns the token
type Token struct {
	JWT  string `json:"jwt"`
	User *User  `json:"user"`
}
