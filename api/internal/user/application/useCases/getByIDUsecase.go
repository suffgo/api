package usecases

import "suffgo/internal/user/domain"

type GetByIDUsecase struct {
	userGetByIDRepository domain.UserRepository
}


func NewGetByIDUsecase(repository domain.UserRepository) *GetByIDUsecase{
	return &GetByIDUsecase{
		userGetByIDRepository: repository,
	}
}