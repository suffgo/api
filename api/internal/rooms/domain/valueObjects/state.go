package valueobjects

import "errors"

type State struct{
	CurrentState string
}

func NewState(value string) (*State, error) {
	states := getStates()
	for _, state := range states {
		if state == value {
			return &State{
				CurrentState: state,
			}, nil
		}
	}

	return nil, errors.ErrUnsupported
}

func (s *State) SetState(value string) error {
	states := getStates()
	for _, state := range states {
		if state == value {
			s.CurrentState = value
		}
	}

	return errors.ErrUnsupported
}

func getStates() [4]string{
	return [...]string{"created", "online", "in progress", "finished"}
}
