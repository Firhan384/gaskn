package dto

type UserAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
