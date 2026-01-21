package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/events"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

type UpdateCustomerInput struct {
	Code     int64
	Customer *domain.Customer
}

type UpdateCustomer struct {
	Repo   ports.CustomerRepository
	Tx     ports.TxManager
	Outbox ports.OutboxWriter
	Topic  string
}

func (uc *UpdateCustomer) Run(ctx context.Context, in UpdateCustomerInput) (*domain.Customer, error) {
	var updated *domain.Customer
	err := uc.Tx.InTx(ctx, func(tx ports.Tx) error {
		before, err := uc.Repo.Get(tx.Context(), in.Code)
		if err != nil {
			return err
		}
		if before == nil {
			return domain.ErrCustomerNotFound
		}
		if in.Customer == nil {
			return domain.ErrCustomerBadInput
		}

		in.Customer.Code = in.Code
		createdAt := before.CreatedAt
		*before = *in.Customer
		before.CreatedAt = createdAt
		before.UpdatedAt = time.Now().UTC()
		if err := uc.Repo.Update(tx.Context(), before); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrCustomerNotFound
			}
			return err
		}

		updated = before
		payload, headers, _ := events.BuildUpdated(uc.Topic, before, "customer-mdm/1.0.0")
		return uc.Outbox.Write(tx, uc.Topic, nil, payload, headers)
	})
	return updated, err
}
