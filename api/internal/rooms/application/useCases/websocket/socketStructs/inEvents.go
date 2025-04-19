package socketStructs

import (
	"encoding/json"
	"errors"
	"log"
	opterr "suffgo/internal/options/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
	votedom "suffgo/internal/votes/domain"
)

// Si el id es = 0, no voto nada
func ReceiveVote(event Event, c *Client) error {

	defer func() {
		c.lobby.votesProcesing <- struct{}{}
	}()

	//si ya voto no podes volver a hacerlo
	if c.voted {
		return nil
	}

	var voteEvent VoteEvent

	if err := json.Unmarshal(event.Payload, &voteEvent); err != nil {
		return err
	}

	votedOpt, err := sv.NewID(uint(voteEvent.OptionId))
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	//chequeo que la opcion exista
	_, err = c.lobby.optRepo.GetByID(*votedOpt)
	if errors.Is(err, opterr.ErrOptNotFound) {
		log.Printf("user %s voted in blank \n", c.User.Username().Username)
	}

	userId := c.User.ID()
	vote := votedom.NewVote(nil, &userId, votedOpt)
	<-c.lobby.votesProcesing
	vote, err = c.lobby.voteRepo.Save(*vote)

	if err != nil {
		log.Println(err.Error()) 
		return nil
	}
	c.lobby.results[c] = *vote

	c.voted = true
	c.lobby.broadcastClientList() //con esto informo el momento en que un usuario vota

	return nil
}

func KickUser(event Event, c *Client) error {
	var kickEvent *KickUserEvent

	if err := json.Unmarshal(event.Payload, &kickEvent); err != nil {
		errorEvent := Event{
			Action:  EventError,
			Payload: marshalOrPanic(ErrorEvent{Message: "unmarshalling error"}),
		}

		c.egress <- errorEvent

		return nil
	}

	if c.lobby.admin.User.ID().Id != c.User.ID().Id {
		errorEvent := Event{
			Action:  EventError,
			Payload: marshalOrPanic(ErrorEvent{Message: "lack of privileges"}),
		}

		c.egress <- errorEvent
		return nil
	}

	clientKicked := false
	for client := range c.lobby.clients {
		if client.User.ID().Id == kickEvent.UserId {
			errorEvent := Event{
				Action:  EventKickUser,
				Payload: marshalOrPanic(ErrorEvent{Message: "you were kicked out of the room"}),
			}

			client.egress <- errorEvent
			<-client.errorSent
			c.lobby.removeClient(client)
			clientKicked = true

		} else {
			errorEvent := Event{
				Action:  EventKickInfoUser,
				Payload: marshalOrPanic(ErrorEvent{Message: "an user was kicked"}),
			}

			client.egress <- errorEvent
		}
	}

	if clientKicked {
		log.Printf("User with id = %d deleted \n", kickEvent.UserId)
	} else {
		log.Println("User to kick not found")
	}

	return nil
}
