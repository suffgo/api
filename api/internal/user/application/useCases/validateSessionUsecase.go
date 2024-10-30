package usecases

import (
	"suffgo/internal/user/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type ValidateSessionUsecase struct{
	repository domain.UserRepository
}

func  NewValidateSessionUsecase(repo domain.UserRepository) *ValidateSessionUsecase {
	return &ValidateSessionUsecase{
		repository: repo,
	}
}

func (v *ValidateSessionUsecase) Execute(session string, id *sv.ID) error {
	
	sessionID, err := v.repository.GetIDBySession(session, *id)

	if err != nil {
		return err
	}

	if sessionID == nil {
		return err
	}

	return nil
}