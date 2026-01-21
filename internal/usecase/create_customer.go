package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/events"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

type CreateCustomerInput struct {
	Customer *domain.Customer
	IdemKey  string
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
		if in.Customer == nil {
			return domain.ErrCustomerBadInput
		}
		now := time.Now().UTC()
		created = in.Customer
		if created.CreatedAt.IsZero() {
			created.CreatedAt = now
		}
		created.UpdatedAt = now
		if err := uc.Repo.Create(tx.Context(), created); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				updated, upsertErr := uc.upsertOnConflict(tx.Context(), created, pgErr)
				if upsertErr != nil {
					return upsertErr
				}
				created = updated
				payload, headers, _ := events.BuildUpdated(uc.Topic, created, "customer-mdm/1.0.0")
				return uc.Outbox.Write(tx, uc.Topic, nil, payload, headers)
			}
			return err
		}

		payload, headers, _ := events.BuildCreated(uc.Topic, created, "customer-mdm/1.0.0")
		return uc.Outbox.Write(tx, uc.Topic, nil, payload, headers)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *CreateCustomer) upsertOnConflict(ctx context.Context, incoming *domain.Customer, pgErr *pgconn.PgError) (*domain.Customer, error) {
	fields := conflictFields(incoming, pgErr)
	if len(fields) == 0 {
		return nil, domain.NewCustomerConflict()
	}
	var existing *domain.Customer
	var err error
	if contains(fields, "email") {
		existing, err = uc.Repo.GetByEmail(ctx, strings.TrimSpace(incoming.Email))
	} else if contains(fields, "phone") {
		existing, err = uc.Repo.GetByPhone(ctx, strings.TrimSpace(incoming.Phone))
	}
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.NewCustomerConflict(fields...)
	}
	now := time.Now().UTC()
	incoming.Code = existing.Code
	incoming.CreatedAt = existing.CreatedAt
	incoming.UpdatedAt = now
	if err := uc.Repo.Update(ctx, incoming); err != nil {
		return nil, err
	}
	return incoming, nil
}

func conflictFields(c *domain.Customer, pgErr *pgconn.PgError) []string {
	if pgErr != nil {
		switch pgErr.ConstraintName {
		case "customers_email_idx":
			return []string{"email"}
		case "customers_phone_idx":
			return []string{"phone"}
		}
	}
	fields := make([]string, 0, 2)
	if strings.TrimSpace(c.Email) != "" {
		fields = append(fields, "email")
	}
	if strings.TrimSpace(c.Phone) != "" {
		fields = append(fields, "phone")
	}
	return fields
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}
