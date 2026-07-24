package repository

import (
	"product-service/internal/entity"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	FindByID(id uint) (*entity.Product, error)
	FindAll() ([]*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll() ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Update(product *entity.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Product{}, id).Error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}
