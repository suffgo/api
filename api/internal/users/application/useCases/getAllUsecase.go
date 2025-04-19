package usecases

import "suffgo/internal/users/domain"

type (
	GetAllUsecase struct {
		getAllRepository domain.UserRepository
	}
)

func NewGetAllUsecase(repository domain.UserRepository) *GetAllUsecase {
	return &GetAllUsecase{
		getAllRepository: repository,
	}
}

func (s *GetAllUsecase) Execute() ([]domain.User, error) {

	users, err := s.getAllRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return users, nil
}
