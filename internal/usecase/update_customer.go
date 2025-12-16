package usecase

import (
	"context"
	"time"

	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/events"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

type UpdateCustomerInput struct {
	ID                                 string
	Name, Email, TaxID, Phone, Country string
	IsActive                           bool
	Attributes                         map[string]string
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
		before, err := uc.Repo.Get(tx.Context(), in.ID)
		if err != nil {
			return err
		}
		if before == nil {
			return nil
		}

		before.Name = in.Name
		before.Email = in.Email
		before.TaxID = in.TaxID
		before.Phone = in.Phone
		before.Country = in.Country
		before.IsActive = in.IsActive
		before.Attributes = in.Attributes
		before.UpdatedAt = time.Now().UTC()
		if err := uc.Repo.Update(tx.Context(), before); err != nil {
			return err
		}

		updated = before
		payload, headers, _ := events.BuildUpdated(uc.Topic, before, "customer-mdm/1.0.0")
		return uc.Outbox.Write(tx, uc.Topic, []byte(before.ID), payload, headers)
	})
	return updated, err
}
