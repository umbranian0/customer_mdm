package usecase

import (
    "context"
    "github.com/umbranian0/customer-mdm/internal/domain"
    "github.com/umbranian0/customer-mdm/internal/ports"
)

type GetCustomer struct { Repo ports.CustomerRepository }

func (uc *GetCustomer) Run(ctx context.Context, id string) (*domain.Customer, error) {
    return uc.Repo.Get(ctx, id)
}
