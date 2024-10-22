package usecases

import "suffgo/internal/user/domain"

type (
	GetAllUsecase struct {
		getAllRepository domain.UserRepository
	}
)

func NewGetAllUsecase(repository domain.UserRepository) *GetAllUsecase{
	return &GetAllUsecase{
		getAllRepository: repository,
	}
}