package socketStructs

import (
	"encoding/json"
	"log"
	sv "suffgo/internal/shared/domain/valueObjects"
	votedom "suffgo/internal/votes/domain"
)

// Si el id es = 0, no voto nada
func ReceiveVote(event Event, c *Client) error {

	defer func() {
		c.lobby.votesProcesing <- struct{}{}
	}()

	var voteEvent VoteEvent

	if err := json.Unmarshal(event.Payload, &voteEvent); err != nil {
		return err
	}

	votedOpt, err := sv.NewID(uint(voteEvent.OptionId))
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	userId := c.user.ID()
	vote := votedom.NewVote(nil, &userId, votedOpt)
	<-c.lobby.votesProcesing
	vote, err = c.lobby.voteRepo.Save(*vote)

	if err != nil {
		log.Println(err.Error()) //TODO: manejar mejor el error en caso de que por alguna razon no se ingrese un id de opt valido
		return nil
	}
	c.lobby.results[*c] = *vote

	return nil
}
