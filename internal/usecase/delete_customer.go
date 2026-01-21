package usecase

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/events"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

type DeleteCustomer struct {
	Repo   ports.CustomerRepository
	Tx     ports.TxManager
	Outbox ports.OutboxWriter
	Topic  string
}

func (uc *DeleteCustomer) Run(ctx context.Context, code int64) error {
	return uc.Tx.InTx(ctx, func(tx ports.Tx) error {
		before, err := uc.Repo.Get(tx.Context(), code)
		if err != nil {
			return err
		}
		if before == nil {
			return domain.ErrCustomerNotFound
		}
		if err := uc.Repo.Delete(tx.Context(), code); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrCustomerNotFound
			}
			return err
		}

		payload, headers, _ := events.BuildDeleted(uc.Topic, code, "customer-mdm/1.0.0")
		return uc.Outbox.Write(tx, uc.Topic, nil, payload, headers)
	})
}
