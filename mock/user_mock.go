package mock

import "go-patron/internal/entity"

type MockRepository struct{}

func (m *MockRepository) Create(user *entity.User) error {
	// Implementación mock para crear usuario
	return nil
}

func (m *MockRepository) FindByID(id uint) (*entity.User, error) {
	return &entity.User{
		ID:   id,
		Name: "Mock User",
	}, nil
}

func (m *MockRepository) Update(user *entity.User) error {
	// Implementación mock para actualizar usuario
	return nil
}

func (m *MockRepository) Delete(id uint) error {
	// Implementación mock para eliminar usuario
	return nil
}
