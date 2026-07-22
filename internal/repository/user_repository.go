package repository

import (
	"go-patron/internal/entity"

	"gorm.io/gorm"
)

// Aquí solo definimos contratos.
type UserRepository interface {
	Create(
		user *entity.User,
	) error

	FindByID(
		id uint,
	) (*entity.User, error)

	Update(
		user *entity.User,
	) error

	Delete(
		id uint,
	) error

	FindAll() []*entity.User
}

// Estructura: Es la implementación concreta de la interfaz.
type userRepository struct {
	db *gorm.DB
}

// Implementación
func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

func (r *userRepository) FindAll() []*entity.User {
	var users []*entity.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil
	}
	return users
}

// Constructor
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

//Devuelve un puntero de la estructura privada (&userRepository),
//pero tipado como la interfaz pública. Esto oculta los detalles de
// la implementación.
