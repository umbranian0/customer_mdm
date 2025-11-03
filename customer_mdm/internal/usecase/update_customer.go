package usecase

import (
    "context"
    "time"

    "github.com/yourorg/customer-mdm/internal/ports"
    customerv1 "github.com/yourorg/customer-mdm/api/gen/customer/v1"
    "google.golang.org/protobuf/types/known/timestamppb"
    "google.golang.org/protobuf/proto"
)

type UpdateCustomerInput struct {
    ID string
    Name, Email, TaxID, Phone, Country string
    IsActive   bool
    Attributes map[string]string
}

type UpdateCustomer struct {
    Repo   ports.CustomerRepository
    Tx     ports.TxManager
    Outbox ports.OutboxWriter
    Topic  string
}

func (uc *UpdateCustomer) Run(ctx context.Context, in UpdateCustomerInput) error {
    return uc.Tx.InTx(ctx, func(tx ports.Tx) error {
        before, err := uc.Repo.Get(tx.Context(), in.ID)
        if err != nil { return err }
        if before == nil { return nil }

        before.Name = in.Name
        before.Email = in.Email
        before.TaxID = in.TaxID
        before.Phone = in.Phone
        before.Country = in.Country
        before.IsActive = in.IsActive
        before.Attributes = in.Attributes
        before.UpdatedAt = time.Now().UTC()
        if err := uc.Repo.Update(tx.Context(), before); err != nil { return err }

        ev := &customerv1.CustomerEvent{
            EventId: in.ID,
            AggregateId: in.ID,
            EventType: "CustomerUpdated",
            OccurredAt: timestamppb.Now(),
            Source: "customer-mdm/1.0.0",
            SchemaVersion: "v1",
            Data: &customerv1.CustomerEvent_Updated{Updated: &customerv1.CustomerUpdated{
                After: &customerv1.Customer{
                    Id: before.ID, Name: before.Name, Email: before.Email, TaxId: before.TaxID,
                    Phone: before.Phone, Country: before.Country, IsActive: before.IsActive,
                    Attributes: before.Attributes,
                    CreatedAt: timestamppb.New(before.CreatedAt), UpdatedAt: timestamppb.New(before.UpdatedAt),
                },
            }},
        }
        payload, _ := proto.Marshal(ev)
        return uc.Outbox.Write(tx, uc.Topic, []byte(before.ID), payload, nil)
    })
}
