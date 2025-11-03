package ports

import (
    "context"
    "github.com/yourorg/customer-mdm/internal/domain"
)

type CustomerRepository interface {
    Create(ctx context.Context, c *domain.Customer) error
    Get(ctx context.Context, id string) (*domain.Customer, error)
    Update(ctx context.Context, c *domain.Customer) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, pageSize int, pageToken, query string) (items []*domain.Customer, nextToken string, total int, err error)
}
