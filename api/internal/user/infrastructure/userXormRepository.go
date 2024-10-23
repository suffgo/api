package infrastructure

import (
	"errors"
	"suffgo/cmd/database"
	d "suffgo/internal/user/domain"
	v "suffgo/internal/user/domain/valueObjects"
	"suffgo/internal/user/infrastructure/mappers"
	m "suffgo/internal/user/infrastructure/models"
)

type UserXormRepository struct {
	db database.Database
}

func NewUserXormRepository(db database.Database) *UserXormRepository {
	return &UserXormRepository{
		db: db,
	}
}

func (s *UserXormRepository) GetByID(id v.ID) (*d.User, error) {
	userModel := new(m.User)
	has, err := s.db.GetDb().ID(id.Id).Get(userModel)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}

	userEnt, err := mappers.ModelToDomain(userModel)

	if err != nil {
		return nil, errors.New("Error de mapeo de datos")
	}

	return userEnt, nil
}

func (s *UserXormRepository) GetAll() ([]d.User, error) {
	var users []d.User
	err := s.db.GetDb().Find(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserXormRepository) Delete(id v.ID) error {

	_, err := s.db.GetDb().ID(id.Id).Delete(&m.User{})
	if err != nil {
		// Manejar el error
		return  err
	}
	return nil
}

func (s *UserXormRepository) Create(user d.User) error {
	return nil
}

func (s *UserXormRepository) GetByEmail(email v.Email) (*d.User, error) {
	return nil, nil
}

func (s *UserXormRepository) Save(user d.User) error {
	userModel := &m.User{
		Dni:      user.Dni().Dni,
		Username: user.Username().Username,
		Password: user.Password().Password,
		Name:     user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Email:    user.Email().Email,
	}

	// Inserta el usuario en la base de datos
	_, err := s.db.GetDb().Insert(userModel)
	if err != nil {
		return err
	}

	return nil
}
