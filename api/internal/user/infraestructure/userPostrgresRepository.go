package infraestructure

import (
	d "suffgo/internal/user/domain"
	v "suffgo/internal/user/domain/valueObjects"
    m "suffgo/internal/user/infraestructure/models"

	"xorm.io/xorm"
)

type UserPostgresRepository struct {
    db *xorm.Engine
}

func NewUserPostgresRepository(db *xorm.Engine) *UserPostgresRepository {
    return &UserPostgresRepository{
        db: db,
    }
}

func (s *UserPostgresRepository) GetByID(id v.UserID) (*d.User, error) {

	return nil, nil
}

func (s *UserPostgresRepository) GetAll() ([]d.User, error) {
	var users []d.User
    err := s.db.Find(&users)
    if err != nil {
        return nil, err
    }
    return users, nil
}

func (s *UserPostgresRepository) Delete(id v.UserID) error {
    return nil
}
func (s *UserPostgresRepository) Create(user d.User) error {
    return nil
}

func (s *UserPostgresRepository) GetByEmail(email v.UserEmail) (*d.User, error) {
    return nil, nil
}

func (s *UserPostgresRepository) Save(user d.User) error {
    userModel := &m.User{
        Dni:      user.Dni().Dni,
        Username: user.Username().Username,
        Password: user.Password().Password,
        Name:     user.FullName().Name,
        Lastname: user.FullName().Lastname,
        Email:    user.Email().Email,
    }

    // Inserta el usuario en la base de datos
    _, err := s.db.Insert(userModel)
    if err != nil {
        return err
    }

    return nil
}
