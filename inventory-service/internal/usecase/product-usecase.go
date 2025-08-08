package usecase

import (
	"context"

	"github.com/Nurda-zh/a1/inventory-service/internal/entity"
	"github.com/Nurda-zh/a1/inventory-service/internal/repository"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, p *entity.Product) error
	GetProduct(ctx context.Context, id string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, id string, p *entity.Product) error
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context) ([]entity.Product, error)
}

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(r repository.ProductRepository) ProductUsecase {
	return &productUsecase{repo: r}
}

func (u *productUsecase) CreateProduct(ctx context.Context, p *entity.Product) error {
	return u.repo.Create(ctx, p)
}

func (u *productUsecase) GetProduct(ctx context.Context, id string) (*entity.Product, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, id string, p *entity.Product) error {
	return u.repo.Update(ctx, id, p)
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *productUsecase) ListProducts(ctx context.Context) ([]entity.Product, error) {
	return u.repo.List(ctx)
}
