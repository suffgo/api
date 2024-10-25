package models

type ProposalOption struct {
	ID         uint `xorm:"'id' pk autoincr"`
	OptionID   uint `xorm:"'option_id' index not null"`
	ProposalID uint `xorm:"'proposal_id' index not null"`
}
