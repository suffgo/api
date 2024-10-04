package models

type AddUserData struct {
	Dni      uint32 `json:"Dni"`
	Mail     string `json:"Mail"`
	Password string `json:"Password"`
	Username string `json:"Username"`
}
