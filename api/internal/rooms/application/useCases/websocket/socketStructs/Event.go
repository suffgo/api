package socketStructs

import "encoding/json"

type Event struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type EventUsecase func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventUpdateClientList = "update_client_list"
	EventStartVoting = "start_voting"
	EventEndVoting = "end_voting"
	EventVote = "vote"
	EventResults = "results"
	EventFirstProp = "first_proposal"
	EventNextProp = "next_proposal"
	EventError = "error"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type UpdateClientListEvent struct {
	Clients []string `json:"clients"`
}

type ErrorEvent struct{
	Message string `json:"message"`
}

type FirstPropEvent struct {
	Prop string `json:"prop"`
}

type NextPropEvent struct {
	Prop string `json:"prop"`
}