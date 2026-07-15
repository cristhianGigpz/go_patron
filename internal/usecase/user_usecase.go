package usecase

import (
	"errors"
	"go-patron/internal/dto"
	"go-patron/internal/entity"
	"go-patron/internal/repository"
)

type UserUseCaseInterface interface {
	Create(req dto.CreateUserRequest) (dto.CreateUserResponse, error)
	FindByID(id uint) (dto.UserResponse, error)
}

type UserUseCase struct {
	repo repository.UserRepository
	//Depende de la interfaz del repositorio y no de una implementación concreta
}

// Observa que el Use Case no sabe si usamos GORM, PostgreSQL o MongoDB.
func (u *UserUseCase) Create(req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	// 1. Validaciones de negocio adicionales
	if req.Name == "" {
		return dto.CreateUserResponse{}, errors.New("nombre requerido")
	}

	if req.Email == "" {
		return dto.CreateUserResponse{}, errors.New("email requerido")
	}
	// 2. [Opcional] Lógica de negocio: Encriptar password

	// 3. Mapear de DTO a Entidad
	userEntity := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Age:      req.Age,
	}

	// 4. Guardar en la base de datos a través del repositorio
	err := u.repo.Create(userEntity)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	// 5. Mapear de Entidad a DTO de respuesta (No devolvemos el Password)
	response := dto.CreateUserResponse{
		Name:  userEntity.Name,
		Email: userEntity.Email,
		Age:   userEntity.Age,
	}

	return response, nil
}

func (u *UserUseCase) FindByID(id uint) (dto.UserResponse, error) {
	if id == 0 {
		return dto.UserResponse{}, errors.New("ID inválido")
	}

	// 1. Obtener la entidad pura desde el repositorio
	userEntity, err := u.repo.FindByID(id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// 2. Mapear/Transformar la entidad al DTO de respuesta seguro
	response := dto.UserResponse{
		ID:    userEntity.ID,
		Name:  userEntity.Name,
		Email: userEntity.Email,
	}

	return response, nil
}

// Constructor
func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

//Recibe el repositorio por medio de inyección de dependencias.
//facilita las pruebas unitarias, ya que permite pasar un mock (falso repositorio) en lugar de una base de datos real.
