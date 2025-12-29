package service

import (
	"context"

	"kenjix.com/persist/internal/model"
	"kenjix.com/persist/internal/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, prod *model.Product) (int64, error)
	GetProduct(ctx context.Context, id int64) (*model.Product, error)
	UpdateProduct(ctx context.Context, prod *model.Product) (int64, error)
	DeleteProduct(ctx context.Context, id int64) (int64, error)
	ListProducts(ctx context.Context) ([]*model.Product, error)
}

// ProductService provides product-related operations
type productServiceImpl struct {
	repo repository.ProductRepository
}

// NewProductService creates a new ProductService
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productServiceImpl{
		repo: repo,
	}
}

// CreateProduct creates a new product
func (s *productServiceImpl) CreateProduct(ctx context.Context, prod *model.Product) (int64, error) {
	return s.repo.Create(ctx, prod)
}

// GetProduct retrieves a product by ID
func (s *productServiceImpl) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateProduct updates an existing product
func (s *productServiceImpl) UpdateProduct(ctx context.Context, prod *model.Product) (int64, error) {
	return s.repo.Update(ctx, prod)
}

// DeleteProduct deletes a product by ID
func (s *productServiceImpl) DeleteProduct(ctx context.Context, id int64) (int64, error) {
	return s.repo.Delete(ctx, id)
}

// ListProducts lists all products
func (s *productServiceImpl) ListProducts(ctx context.Context) ([]*model.Product, error) {
	return s.repo.List(ctx)
}
