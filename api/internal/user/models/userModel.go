package models

type AddUserData struct {
	Dni      string `json:"dni"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
	Username string `json:"username"`
}
