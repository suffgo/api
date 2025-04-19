package usecases

import (
	"suffgo/internal/users/domain"
	v "suffgo/internal/users/domain/valueObjects"
)

type GetByEmailUsecase struct {
	userGetByEmailRepository domain.UserRepository
}

func NewGetByEmailUsecase(repository domain.UserRepository) *GetByEmailUsecase {
	return &GetByEmailUsecase{
		userGetByEmailRepository: repository,
	}
}

func (s *GetByEmailUsecase) Execute(email v.Email) (*domain.User, error) {

	user, err := s.userGetByEmailRepository.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	return user, nil
}
