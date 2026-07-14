package usecase

import (
	"errors"
	"go-patron/internal/entity"
	"go-patron/internal/repository"
)

type UserUseCaseInterface interface {
	Create(user *entity.User) error
	FindByID(id uint) (*entity.User, error)
}

type UserUseCase struct {
	repo repository.UserRepository
}

// Observa que el Use Case no sabe si usamos GORM, PostgreSQL o MongoDB.
func (u *UserUseCase) Create(user *entity.User) error {

	if user.Name == "" {
		return errors.New("nombre requerido")
	}

	return u.repo.Create(user)
}

func (u *UserUseCase) FindByID(id uint) (*entity.User, error) {
	if id == 0 {
		return nil, errors.New("ID inválido")
	}
	return u.repo.FindByID(id)
}

// Constructor
func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}
