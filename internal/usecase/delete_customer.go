package usecase

import (
	"context"

	"github.com/umbranian0/customer-mdm/internal/events"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

type DeleteCustomer struct {
	Repo   ports.CustomerRepository
	Tx     ports.TxManager
	Outbox ports.OutboxWriter
	Topic  string
}

func (uc *DeleteCustomer) Run(ctx context.Context, id string) error {
	return uc.Tx.InTx(ctx, func(tx ports.Tx) error {
		before, err := uc.Repo.Get(tx.Context(), id)
		if err != nil {
			return err
		}
		if before == nil {
			return nil
		}
		if err := uc.Repo.Delete(tx.Context(), id); err != nil {
			return err
		}

		payload, headers, _ := events.BuildDeleted(uc.Topic, id, "customer-mdm/1.0.0")
		return uc.Outbox.Write(tx, uc.Topic, []byte(id), payload, headers)
	})
}
