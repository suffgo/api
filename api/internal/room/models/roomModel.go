package models

type AddRoomData struct {
	LinkInvite *string `json:"linkInvite"`
	IsFormal   bool    `json:"isFormal"`
	Name       string  `json:"name"`
	AdminID    uint    `json:"adminID"`
}
