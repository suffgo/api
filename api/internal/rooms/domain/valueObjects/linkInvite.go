package valueobjects

import "errors"

type (
	LinkInvite struct {
		LinkInvite string
	}
)

func NewLinkInvite(linkInvite string) (*LinkInvite, error) {
	if linkInvite == "" {
		return nil, errors.New("invalid linkInvite")
	}

	return &LinkInvite{
		LinkInvite: linkInvite,
	}, nil
}
