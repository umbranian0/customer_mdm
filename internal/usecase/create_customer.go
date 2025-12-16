package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/events"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

type CreateCustomerInput struct {
	Name, Email, TaxID, Phone, Country string
	IsActive                           bool
	Attributes                         map[string]string
	IdemKey                            string
}

type CreateCustomer struct {
	Repo   ports.CustomerRepository
	Tx     ports.TxManager
	Outbox ports.OutboxWriter
	Topic  string
}

func (uc *CreateCustomer) Run(ctx context.Context, in CreateCustomerInput) (*domain.Customer, error) {
	var created *domain.Customer
	err := uc.Tx.InTx(ctx, func(tx ports.Tx) error {
		now := time.Now().UTC()
		created = &domain.Customer{
			ID:   uuid.New().String(),
			Name: in.Name, Email: in.Email, TaxID: in.TaxID, Phone: in.Phone,
			Country: in.Country, IsActive: in.IsActive, Attributes: in.Attributes,
			CreatedAt: now, UpdatedAt: now,
		}
		if err := uc.Repo.Create(tx.Context(), created); err != nil {
			return err
		}

		payload, headers, _ := events.BuildCreated(uc.Topic, created, "customer-mdm/1.0.0")
		return uc.Outbox.Write(tx, uc.Topic, []byte(created.ID), payload, headers)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
