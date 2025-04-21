package models

type Proposal struct {
	ID          uint    `xorm:"'id' pk autoincr"`
	Archive     *string `xorm:"'archive' null"` // Archivo con informacion detallada de la propuesta
	Title       string  `xorm:"'title' not null"`
	Description *string `xorm:"'description' null"`
	RoomID      uint    `xorm:"'room_id' index not null"`
}

type SqlResult struct {
	ProposalId          uint   `xorm:"proposal_id"`
	ProposalTitle       string `xorm:"proposal_title"`
	ProposalDescription string `xorm:"proposal_description"`
	RoomID              uint   `xorm:"room_id"`
	OptionId            uint   `xorm:"option_id"`
	OptionValue         string `xorm:"option_value"`
	VoteId              uint   `xorm:"vote_id"`
	UserId              uint   `xorm:"user_id"`
	Username            string `xorm:"username"`
	UserImage           string `xorm:"user_image"`
}
