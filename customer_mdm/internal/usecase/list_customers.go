package usecase

import (
    "context"
    "github.com/yourorg/customer-mdm/internal/domain"
    "github.com/yourorg/customer-mdm/internal/ports"
)

type ListCustomersInput struct {
    PageSize int
    PageToken string
    Query string
}

type ListCustomers struct { Repo ports.CustomerRepository }

func (uc *ListCustomers) Run(ctx context.Context, in ListCustomersInput) ([]*domain.Customer, string, int, error) {
    return uc.Repo.List(ctx, in.PageSize, in.PageToken, in.Query)
}
