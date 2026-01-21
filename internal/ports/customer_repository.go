package ports

import (
	"context"

	"github.com/umbranian0/customer-mdm/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, c *domain.Customer) error
	Get(ctx context.Context, code int64) (*domain.Customer, error)
	GetByEmail(ctx context.Context, email string) (*domain.Customer, error)
	GetByPhone(ctx context.Context, phone string) (*domain.Customer, error)
	Update(ctx context.Context, c *domain.Customer) error
	Delete(ctx context.Context, code int64) error
	List(ctx context.Context, pageSize int, pageToken, query string) (items []*domain.Customer, nextToken string, total int, err error)
}
