package socketStructs

import (
	"log"
	optdom "suffgo/internal/options/domain"
)

func StartVoting(event Event, c *Client) error {
	log.Printf("room with id = %d has begun \n", c.Lobby().room.AdminID().Id)

	if c.user.ID().Id != c.Lobby().Admin().user.ID().Id {

		errorEvent := Event{
			Action:  EventError,
			Payload: marshalOrPanic(ErrorEvent{Message: "You are not the admin"}),
		}

		c.egress <- errorEvent
		return nil
	}

	//esto deberia ser chequeado antes, no deberia poder comenzar una sala que no tiene propuestas
	if len(c.Lobby().proposals) >= c.lobby.nextProposal {

		lastProp := false
		if c.lobby.nextProposal == len(c.Lobby().proposals)-1 {
			lastProp = true
		}

		proposal := c.Lobby().proposals[c.lobby.nextProposal]
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
			LastProp:    lastProp,
		}

		prop := Event{
			Action:  EventFirstProp,
			Payload: marshalOrPanic(proposalevt),
		}

		for client := range c.Lobby().Clients() {
			client.egress <- prop
		}
	}

	c.lobby.room.State().SetState("in progress")
	_, err := c.lobby.roomRepo.Update(c.lobby.room)

	if err != nil {
		return nil
	}

	log.Println(c.lobby.room.State().CurrentState)
	c.lobby.nextProposal++

	return nil
}

func NextProposal(event Event, c *Client) error {

	log.Println("enviando next proposal")
	if c.user.ID().Id != c.Lobby().Admin().user.ID().Id {

		errorEvent := Event{
			Action:  EventError,
			Payload: marshalOrPanic(ErrorEvent{Message: "You are not the admin"}),
		}

		c.egress <- errorEvent
		return nil
	}

	lastProp := false
	if c.lobby.nextProposal == len(c.Lobby().proposals)-1 {
		lastProp = true
	}

	if c.lobby.nextProposal <= len(c.Lobby().proposals) {
		proposal := c.Lobby().proposals[c.lobby.nextProposal]
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
			LastProp:    lastProp,
		}

		prop := Event{
			Action:  EventNextProp,
			Payload: marshalOrPanic(proposalevt),
		}

		for client := range c.Lobby().Clients() {
			client.voted = false
			client.egress <- prop
		}
	} else {
		log.Println("no more proposals")
		return nil
	}

	c.lobby.broadcastClientList()
	c.lobby.nextProposal++
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

func SendResults(event Event, c *Client) error {
	//armo el json con los votos
	var userVotes []UserVoteEvent
	for client, vote := range c.Lobby().results {
		voterData := VoterData{
			Username: client.user.Username().Username,
			ID:       client.user.ID().Id,
		}

		userVote := UserVoteEvent{
			From:     voterData,
			OptionId: vote.OptionID().Id,
		}

		userVotes = append(userVotes, userVote)
	}

	evt := Event{
		Action:  EventResults,
		Payload: marshalOrPanic(userVotes),
	}

	log.Println(evt)
	for client := range c.Lobby().Clients() {
		client.egress <- evt
		log.Println("resultados enviados a " + client.user.Username().Username)
	}

	return nil
}
