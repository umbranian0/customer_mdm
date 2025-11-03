package usecase

import (
	"context"

	customerv1 "github.com/umbranian0/customer-mdm/api/gen/customer/v1"
	"github.com/umbranian0/customer-mdm/internal/ports"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

		ev := &customerv1.CustomerEvent{
			EventId:       id,
			AggregateId:   id,
			EventType:     "CustomerDeleted",
			OccurredAt:    timestamppb.Now(),
			Source:        "customer-mdm/1.0.0",
			SchemaVersion: "v1",
			Data:          &customerv1.CustomerEvent_Deleted{Deleted: &customerv1.CustomerDeleted{}}, // minimal
		}
		payload, _ := proto.Marshal(ev)
		return uc.Outbox.Write(tx, uc.Topic, []byte(id), payload, nil)
	})
}
