package repository

import (
	"context"

	"kenjix.com/persist/internal/model"
)

// ProductRepository defines behavior for product persistence.
type ProductRepository interface {
	Create(ctx context.Context, prod *model.Product) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Product, error)
	Update(ctx context.Context, prod *model.Product) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
	List(ctx context.Context) ([]*model.Product, error)
}
