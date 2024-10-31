package infrastructure

import (
	"suffgo/cmd/database"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
	d "suffgo/internal/users/domain"
	ue "suffgo/internal/users/domain/errors"
	v "suffgo/internal/users/domain/valueObjects"
	"suffgo/internal/users/infrastructure/mappers"
	m "suffgo/internal/users/infrastructure/models"
)

type UserXormRepository struct {
	db database.Database
}

func NewUserXormRepository(db database.Database) *UserXormRepository {
	return &UserXormRepository{
		db: db,
	}
}

func (s *UserXormRepository) GetByID(id sv.ID) (*d.User, error) {
	userModel := new(m.Users)
	has, err := s.db.GetDb().ID(id.Id).Get(userModel)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ue.UserNotFoundError
	}

	userEnt, err := mappers.ModelToDomain(userModel)

	if err != nil {
		return nil, se.DataMappingError
	}

	return userEnt, nil
}

func (s *UserXormRepository) GetAll() ([]d.User, error) {
	var users []m.Users
	err := s.db.GetDb().Find(&users)
	if err != nil {
		return nil, err
	}

	var usersDomain []d.User
	for _, user := range users {
		userDomain, err := mappers.ModelToDomain(&user)

		if err != nil {
			return nil, err
		}

		usersDomain = append(usersDomain, *userDomain)
	}
	return usersDomain, nil
}

func (s *UserXormRepository) Delete(id sv.ID) error {

	affected, err := s.db.GetDb().ID(id.Id).Delete(&m.Users{})
	if err != nil {
		return err
	}

	if affected == 0 {
		return ue.UserNotFoundError
	}

	return nil
}

func (s *UserXormRepository) GetByEmail(email v.Email) (*d.User, error) {
	var user m.Users
	has, err := s.db.GetDb().Where("email = ?", email.Email).Get(&user)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	mappedUser, err := mappers.ModelToDomain(&user)

	if err != nil {
		return nil, err
	}
	return mappedUser, nil
}

func (s *UserXormRepository) GetByDni(dni v.Dni) (*d.User, error) {
	var user m.Users
	has, err := s.db.GetDb().Where("dni = ?", dni.Dni).Get(&user)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	mappedUser, err := mappers.ModelToDomain(&user)

	if err != nil {
		return nil, err
	}
	return mappedUser, nil
}

func (s *UserXormRepository) GetByUsername(username v.UserName) (*d.User, error) {
	var user m.Users
	has, err := s.db.GetDb().Where("username = ?", username.Username).Get(&user)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	mappedUser, err := mappers.ModelToDomain(&user)

	if err != nil {
		return nil, err
	}
	return mappedUser, nil
}

func (s *UserXormRepository) Save(user d.User) error {
	userModel := &m.Users{
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
