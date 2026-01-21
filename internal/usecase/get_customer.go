package usecase

import (
	"context"

	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

type GetCustomer struct{ Repo ports.CustomerRepository }

func (uc *GetCustomer) Run(ctx context.Context, code int64) (*domain.Customer, error) {
	c, err := uc.Repo.Get(ctx, code)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, domain.ErrCustomerNotFound
	}
	return c, nil
}
