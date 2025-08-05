package models

type UserProfile struct {
	ID       uint    `json:"id" example:"1"`
	Username string  `json:"username" example:"admin"`
	Role     string  `json:"role" example:"admin"`
	Balance  float64 `json:"balance" example:"0.0"`
}
