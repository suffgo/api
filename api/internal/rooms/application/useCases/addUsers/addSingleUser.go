package addusers

import "suffgo/internal/rooms/domain"

type AddSingleUserUsecase struct {
	repository domain.RoomRepository
}

func NewAddSingleUserUsecase(repository domain.RoomRepository) *AddSingleUserUsecase {
	return &AddSingleUserUsecase{
		repository: repository,
	}
}

func (s *AddSingleUserUsecase) Execute(userData string) error {
	//Tengo que ver que tipo de dato es
	
	//Si es mail
	
	//Obtengo user por mail

	//Si nombre de usuario

	//Obtengo user por nombre de usuario

	//Si es DNI
	//obtengo user por dni

	//Si no existe, error

	//Si existe, lo agrego a la sala

	return nil
}
