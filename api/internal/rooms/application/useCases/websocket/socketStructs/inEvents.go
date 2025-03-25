package socketStructs

import (
	"encoding/json"
	"errors"
	"log"
	sv "suffgo/internal/shared/domain/valueObjects"
	votedom "suffgo/internal/votes/domain"
	opterr "suffgo/internal/options/domain/errors"
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
		log.Println(err.Error()) //TODO: manejar mejor el error en caso de que por alguna razon no se ingrese un id de opt valido
		return nil
	}
	c.lobby.results[c] = *vote

	c.voted = true
	c.lobby.broadcastClientList() //con esto informo el momento en que un usuario vota

	return nil
}
