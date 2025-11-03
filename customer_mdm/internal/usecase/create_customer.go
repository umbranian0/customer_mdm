package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	customerv1 "github.com/umbranian0/customer-mdm/api/gen/customer/v1"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/ports"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

		ev := &customerv1.CustomerEvent{
			EventId:       uuid.New().String(),
			AggregateId:   created.ID,
			EventType:     "CustomerCreated",
			OccurredAt:    timestamppb.New(now),
			Source:        "customer-mdm/1.0.0",
			SchemaVersion: "v1",
			Data: &customerv1.CustomerEvent_Created{Created: &customerv1.CustomerCreated{After: &customerv1.Customer{
				Id:   created.ID,
				Name: created.Name, Email: created.Email, TaxId: created.TaxID, Phone: created.Phone,
				Country: created.Country, IsActive: created.IsActive, Attributes: created.Attributes,
				CreatedAt: timestamppb.New(created.CreatedAt), UpdatedAt: timestamppb.New(created.UpdatedAt),
			}}},
		}
		payload, _ := proto.Marshal(ev)
		return uc.Outbox.Write(tx, uc.Topic, []byte(created.ID), payload, nil)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
