package usecase

import (
	"errors"
	"product-service/internal/dto"
	"product-service/internal/entity"
	"product-service/internal/repository"
)

type ProductUseCaseInterface interface {
	Create(req dto.CreateProductRequest) (dto.ProductResponse, error)
	FindByID(id uint) (dto.ProductResponse, error)
	FindAll() ([]dto.ProductResponse, error)
	Update(id uint, req dto.UpdateProductRequest) (dto.ProductResponse, error)
	Delete(id uint) error
}

type ProductUseCase struct {
	repo repository.ProductRepository
}

func (u *ProductUseCase) Create(req dto.CreateProductRequest) (dto.ProductResponse, error) {
	if req.Name == "" {
		return dto.ProductResponse{}, errors.New("nombre requerido")
	}

	product := &entity.Product{Name: req.Name, Description: req.Description, Price: req.Price, Stock: req.Stock}
	if err := u.repo.Create(product); err != nil {
		return dto.ProductResponse{}, err
	}
	return toResponse(product), nil
}

func (u *ProductUseCase) FindByID(id uint) (dto.ProductResponse, error) {
	if id == 0 {
		return dto.ProductResponse{}, errors.New("ID inválido")
	}
	product, err := u.repo.FindByID(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	return toResponse(product), nil
}

func (u *ProductUseCase) FindAll() ([]dto.ProductResponse, error) {
	products, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		if product != nil {
			responses = append(responses, toResponse(product))
		}
	}
	return responses, nil
}

func (u *ProductUseCase) Update(id uint, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	if id == 0 {
		return dto.ProductResponse{}, errors.New("ID inválido")
	}
	product, err := u.repo.FindByID(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	product.Name, product.Description = req.Name, req.Description
	product.Price, product.Stock = req.Price, req.Stock
	if err := u.repo.Update(product); err != nil {
		return dto.ProductResponse{}, err
	}
	return toResponse(product), nil
}

func (u *ProductUseCase) Delete(id uint) error {
	if id == 0 {
		return errors.New("ID inválido")
	}
	return u.repo.Delete(id)
}

func toResponse(product *entity.Product) dto.ProductResponse {
	return dto.ProductResponse{ID: product.ID, Name: product.Name, Description: product.Description, Price: product.Price, Stock: product.Stock}
}

func NewProductUseCase(repo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}
