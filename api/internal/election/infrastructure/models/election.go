package models

type Election struct {
	ID               uint `xorm:"'id' pk autoincr"`
	ProposalOptionID uint `xorm:"'proposal_option_id' index not null"`
	UserID           uint `xorm:"'user_id' index not null"`
	ProposalID       uint `xorm:"'proposal_id' index not null"`
}