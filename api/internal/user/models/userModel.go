package models

type AddUserData struct {
	Dni      string `json:"Dni"`
	Mail     string `json:"Mail"`
	Password string `json:"Password"`
	Username string `json:"Username"`
}
