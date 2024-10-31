package models

type Option struct {
	ID         uint   `xorm:"'id' pk autoincr"`
	Value      string `xorm:"'value' not null unique"`
	ProposalID uint   `xorm:"'proposal_id' index not null"`
}
