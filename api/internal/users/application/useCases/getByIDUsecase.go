package usecases

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	"suffgo/internal/users/domain"
)

type GetByIDUsecase struct {
	userGetByIDRepository domain.UserRepository
}

func NewGetByIDUsecase(repository domain.UserRepository) *GetByIDUsecase {
	return &GetByIDUsecase{
		userGetByIDRepository: repository,
	}
}

func (s *GetByIDUsecase) Execute(id sv.ID) (*domain.User, error) {

	user, err := s.userGetByIDRepository.GetByID(id)

	if err != nil {
		return nil, err
	}
	return user, nil
}
