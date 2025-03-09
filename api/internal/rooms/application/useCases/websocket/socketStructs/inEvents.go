package socketStructs

import (
	"encoding/json"
	"log"
	sv "suffgo/internal/shared/domain/valueObjects"
	votedom "suffgo/internal/votes/domain"
)

// Si el id es = 0, no voto nada
func ReceiveVote(event Event, c *Client) error {

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
	//cuidado con esto capaz tenga que usar un canal
	_, err = c.lobby.voteRepo.Save(*vote)

	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return nil
}
