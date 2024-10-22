package usecases

import (
	"suffgo/internal/user/domain"
	v "suffgo/internal/user/domain/valueObjects"
)

type GetByIDUsecase struct {
	userGetByIDRepository domain.UserRepository
}


func NewGetByIDUsecase(repository domain.UserRepository) *GetByIDUsecase{
	return &GetByIDUsecase{
		userGetByIDRepository: repository,
	}
}

func (s *GetByIDUsecase) Execute(id v.ID) (*domain.User, error) {

	user, err := s.userGetByIDRepository.GetByID(id)

	if err != nil {
		return nil, err
	}
	return user, nil
}