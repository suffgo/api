package socketStructs

import (
	"encoding/json"
	optdom "suffgo/internal/options/domain"
)

type Event struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type EventUsecase func(event Event, c *Client) error

const (
	EventSendMessage      = "send_message"
	EventUpdateClientList = "update_client_list"
	EventStartVoting      = "start_voting"
	EventEndVoting        = "end_voting"
	EventVote             = "vote"
	EventResults          = "results"
	EventFirstProp        = "first_proposal"
	EventNextProp         = "next_proposal"
	EventError            = "error"
	EventKickUser         = "kick_user"
	EventKickInfoUser     = "kick_info"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type UpdateClientListEvent struct {
	Clients []ClientData `json:"clients"`
}

// estructura utilizada para trackear votos en tiempo real
type ClientData struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Voted    bool   `json:"voted"`
	Image    string `json:"image"`
}

type ErrorEvent struct {
	Message string `json:"message"`
}

type ProposalEvent struct {
	ID          uint               `json:"id"`
	Archive     *string            `json:"archive"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	RoomID      uint               `json:"room_id"`
	Options     []optdom.OptionDTO `json:"options"`
	LastProp    bool               `json:"last_prop"`
}

type NextPropEvent struct {
	Prop string `json:"prop"`
}

type VoteEvent struct {
	OptionId uint `json:"option_id"`
}

type UserVoteEvent struct {
	From     VoterData `json:"from"`
	OptionId uint      `json:"option_id"`
}

type VoterData struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
}

type ResultsEvent struct {
	Votes []UserVoteEvent `json:"votes"`
}

type KickUserEvent struct {
	UserId uint `json:"user_id"`
}
