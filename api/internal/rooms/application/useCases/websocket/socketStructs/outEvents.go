package socketStructs

import (
	"log"
	optdom "suffgo/internal/options/domain"
)

func StartVoting(event Event, c *Client) error {
	log.Printf("room with id = %d has begun \n", c.Lobby().room.AdminID().Id)

	for _, prop := range c.Lobby().proposals {
		log.Println(prop)
	}

	if c.user.ID().Id != c.Lobby().Admin().user.ID().Id {

		errorEvent := Event{
			Action:  EventError,
			Payload: marshalOrPanic(ErrorEvent{Message: "You are not the admin"}),
		}

		c.egress <- errorEvent
		return nil
	}

	//esto deberia ser chequeado antes, no deberia poder comenzar una sala que no tiene propuestas
	if len(c.Lobby().proposals) > 0 {
		proposal := c.Lobby().proposals[0]
		options, err := c.lobby.optRepo.GetByProposal(proposal.ID())
		if err != nil {
			errorEvent := Event{
				Action:  EventError,
				Payload: marshalOrPanic(ErrorEvent{Message: "error fetching options"}),
			}

			c.egress <- errorEvent

			return nil
		}

		var optionsValue []optdom.OptionDTO
		for _, option := range options {

			opt := optdom.OptionDTO{
				ID:         option.ID().Id,
				Value:      option.Value().Value,
				ProposalID: option.ProposalID().Id,
			}
			optionsValue = append(optionsValue, opt)
		}

		proposalevt := ProposalEvent{
			ID:          proposal.ID().Id,
			Archive:     &proposal.Archive().Archive,
			Description: &proposal.Description().Description,
			Title:       proposal.Title().Title,
			Options:     optionsValue,
		}

		prop := Event{
			Action:  EventFirstProp,
			Payload: marshalOrPanic(proposalevt),
		}

		for client := range c.Lobby().Clients() {
			client.egress <- prop
		}
	}

	return nil
}


func SendMessage(event Event, c *Client) error {
	for client := range c.Lobby().Clients() {
		if client != c {
			client.egress <- event
		}
	}
	return nil
}