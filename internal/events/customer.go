package events

import (
	"strings"

	"github.com/google/uuid"
	customerv1 "github.com/umbranian0/customer-mdm/api/gen/customer/v1"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Common headers we add to every customer event for routing/observability.
func baseHeaders(evType, topic string, c *domain.Customer) map[string]string {
	h := map[string]string{
		"event_type":   evType,
		"content_type": "application/json",
		"topic":        topic,
	}
	if c != nil {
		if c.Country != "" {
			h["customer_country"] = c.Country
		}
		h["customer_is_active"] = boolToStr(c.IsActive)
		if dom := emailDomain(c.Email); dom != "" {
			h["customer_email_domain"] = dom
		}
	}
	return h
}

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func emailDomain(email string) string {
	if i := strings.LastIndex(email, "@"); i > 0 && i < len(email)-1 {
		return email[i+1:]
	}
	return ""
}

// BuildCreated builds payload and headers for CustomerCreated.
func BuildCreated(topic string, c *domain.Customer, source string) ([]byte, map[string]string, error) {
	ev := &customerv1.CustomerEvent{
		EventId:       uuid.New().String(),
		AggregateId:   c.ID,
		EventType:     "CustomerCreated",
		OccurredAt:    timestamppb.New(c.CreatedAt),
		Source:        source,
		SchemaVersion: "v1",
		Data: &customerv1.CustomerEvent_Created{Created: &customerv1.CustomerCreated{
			After: toProtoCustomer(c),
		}},
	}
	payload, err := toJSON(ev)
	return payload, baseHeaders(ev.EventType, topic, c), err
}

// BuildUpdated builds payload and headers for CustomerUpdated.
func BuildUpdated(topic string, c *domain.Customer, source string) ([]byte, map[string]string, error) {
	ev := &customerv1.CustomerEvent{
		EventId:       c.ID,
		AggregateId:   c.ID,
		EventType:     "CustomerUpdated",
		OccurredAt:    timestamppb.New(c.UpdatedAt),
		Source:        source,
		SchemaVersion: "v1",
		Data: &customerv1.CustomerEvent_Updated{Updated: &customerv1.CustomerUpdated{
			After: toProtoCustomer(c),
		}},
	}
	payload, err := toJSON(ev)
	return payload, baseHeaders(ev.EventType, topic, c), err
}

// BuildDeleted builds payload and headers for CustomerDeleted.
func BuildDeleted(topic, id, source string) ([]byte, map[string]string, error) {
	ev := &customerv1.CustomerEvent{
		EventId:       id,
		AggregateId:   id,
		EventType:     "CustomerDeleted",
		OccurredAt:    timestamppb.Now(),
		Source:        source,
		SchemaVersion: "v1",
		Data:          &customerv1.CustomerEvent_Deleted{Deleted: &customerv1.CustomerDeleted{}},
	}
	payload, err := toJSON(ev)
	return payload, baseHeaders(ev.EventType, topic, nil), err
}

func toJSON(ev *customerv1.CustomerEvent) ([]byte, error) {
	opts := protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}
	return opts.Marshal(ev)
}

func toProtoCustomer(c *domain.Customer) *customerv1.Customer {
	if c == nil {
		return nil
	}
	return &customerv1.Customer{
		Id:         c.ID,
		Name:       c.Name,
		Email:      c.Email,
		TaxId:      c.TaxID,
		Phone:      c.Phone,
		Country:    c.Country,
		IsActive:   c.IsActive,
		Attributes: c.Attributes,
		CreatedAt:  timestamppb.New(c.CreatedAt),
		UpdatedAt:  timestamppb.New(c.UpdatedAt),
	}
}
