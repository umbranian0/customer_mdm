package events

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	customerv1 "github.com/umbranian0/customer-mdm/api/gen/api/proto/customer/v1"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Common headers we add to every customer event for routing/observability.
func baseHeaders(evType, topic, eventID string, occurredAt *timestamppb.Timestamp, c *domain.Customer) map[string]string {
	h := map[string]string{
		"event_type":   evType,
		"content_type": "application/json",
		"topic":        topic,
		"SN":           "customer",
		"SI":           "upsert",
	}
	if eventID != "" {
		h["ID"] = eventID
	}
	if occurredAt != nil {
		h["created_at"] = occurredAt.AsTime().UTC().Format(time.RFC3339Nano)
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

func AggregateKey(code int64) string {
	if code == 0 {
		return ""
	}
	return strconv.FormatInt(code, 10)
}

// BuildCreated builds payload and headers for CustomerCreated.
func BuildCreated(topic string, c *domain.Customer, source string) ([]byte, map[string]string, error) {
	ev := &customerv1.CustomerEvent{
		EventId:       uuid.New().String(),
		AggregateId:   AggregateKey(c.Code),
		EventType:     "CustomerCreated",
		OccurredAt:    toTimestamp(c.CreatedAt),
		Source:        source,
		SchemaVersion: "v1",
		Data: &customerv1.CustomerEvent_Created{Created: &customerv1.CustomerCreated{
			After: toProtoCustomer(c),
		}},
	}
	payload, err := toJSON(ev)
	return payload, baseHeaders(ev.EventType, topic, ev.EventId, ev.OccurredAt, c), err
}

// BuildUpdated builds payload and headers for CustomerUpdated.
func BuildUpdated(topic string, c *domain.Customer, source string) ([]byte, map[string]string, error) {
	ev := &customerv1.CustomerEvent{
		EventId:       uuid.New().String(),
		AggregateId:   AggregateKey(c.Code),
		EventType:     "CustomerUpdated",
		OccurredAt:    toTimestamp(c.UpdatedAt),
		Source:        source,
		SchemaVersion: "v1",
		Data: &customerv1.CustomerEvent_Updated{Updated: &customerv1.CustomerUpdated{
			After: toProtoCustomer(c),
		}},
	}
	payload, err := toJSON(ev)
	return payload, baseHeaders(ev.EventType, topic, ev.EventId, ev.OccurredAt, c), err
}

// BuildDeleted builds payload and headers for CustomerDeleted.
func BuildDeleted(topic string, code int64, source string) ([]byte, map[string]string, error) {
	ev := &customerv1.CustomerEvent{
		EventId:       uuid.New().String(),
		AggregateId:   AggregateKey(code),
		EventType:     "CustomerDeleted",
		OccurredAt:    timestamppb.Now(),
		Source:        source,
		SchemaVersion: "v1",
		Data:          &customerv1.CustomerEvent_Deleted{Deleted: &customerv1.CustomerDeleted{}},
	}
	payload, err := toJSON(ev)
	return payload, baseHeaders(ev.EventType, topic, ev.EventId, ev.OccurredAt, nil), err
}

func toJSON(ev *customerv1.CustomerEvent) ([]byte, error) {
	opts := protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}
	return opts.Marshal(ev)
}

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func toTimestampPtr(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}

func toProtoCustomer(c *domain.Customer) *customerv1.Customer {
	if c == nil {
		return nil
	}
	return &customerv1.Customer{
		Code:                      c.Code,
		FrontendId:                c.FrontendID,
		ErpId:                     c.ErpID,
		MarketId:                  c.MarketID,
		MarketCustomizerId:        c.MarketCustomizerID,
		Level:                     c.Level,
		ParentId:                  c.ParentID,
		DiscountProfile:           c.DiscountProfile,
		IsActive:                  c.IsActive,
		CanOrder:                  c.CanOrder,
		Username:                  c.Username,
		Password:                  c.Password,
		Email:                     c.Email,
		EmailCopy:                 c.EmailCopy,
		CountryCode:               c.CountryCode,
		Language:                  c.Language,
		ContactLanguage:           c.ContactLanguage,
		WebserviceKey:             c.WebserviceKey,
		Name:                      c.Name,
		Company:                   c.Company,
		TaxId:                     c.TaxID,
		Bank:                      c.Bank,
		BankAddress:               c.BankAddress,
		BankBranch:                c.BankBranch,
		Website:                   c.Website,
		AddressLine1:              c.AddressLine1,
		AddressLine2:              c.AddressLine2,
		PostalCode:                c.PostalCode,
		City:                      c.City,
		Phone:                     c.Phone,
		AccountManagerName:        c.AccountManagerName,
		AccountManagerPhone:       c.AccountManagerPhone,
		AccountManagerEmail:       c.AccountManagerEmail,
		BirthDate:                 toTimestampPtr(c.BirthDate),
		RegisteredAt:              toTimestampPtr(c.RegisteredAt),
		LastLoginAt:               toTimestampPtr(c.LastLoginAt),
		FavoritesNotifications:    c.FavoritesNotifications,
		KeyCode:                   c.KeyCode,
		IsConfirmed:               c.IsConfirmed,
		RecoveryTimestamp:         toTimestampPtr(c.RecoveryTimestamp),
		ReceivesNewsletters:       c.ReceivesNewsletters,
		StandardTier:              c.StandardTier,
		OwnerId:                   c.OwnerID,
		StockPolicy:               c.StockPolicy,
		Location:                  c.Location,
		StreetType:                c.StreetType,
		Neighborhood:              c.Neighborhood,
		State:                     c.State,
		StateRegistration:         c.StateRegistration,
		Country:                   c.Country,
		Comment:                   c.Comment,
		RegistrationCertificate:   c.RegistrationCertificate,
		LastLoginIp:               c.LastLoginIP,
		LastLoginCountryCode:      c.LastLoginCountryCode,
		BlockedBySuspiciousChange: c.BlockedBySuspiciousChange,
		WarehouseCode:             c.WarehouseCode,
		OldErpId:                  c.OldErpID,
		CommercialMarketId:        c.CommercialMarketID,
		Migrated:                  c.Migrated,
		CommercialAreaId:          c.CommercialAreaID,
		IndustrialProduction:      c.IndustrialProduction,
		DeliveryNote:              c.DeliveryNote,
		NoDirectApprovals:         c.NoDirectApprovals,
		IsCleaned:                 c.IsCleaned,
		Addresses:                 toProtoAddresses(c.Addresses),
		CreatedAt:                 toTimestamp(c.CreatedAt),
		UpdatedAt:                 toTimestamp(c.UpdatedAt),
	}
}

func toProtoAddresses(addresses []domain.CustomerAddress) []*customerv1.CustomerAddress {
	if len(addresses) == 0 {
		return nil
	}
	out := make([]*customerv1.CustomerAddress, 0, len(addresses))
	for _, a := range addresses {
		out = append(out, &customerv1.CustomerAddress{
			Id:            a.ID,
			CustomerCode:  a.CustomerCode,
			ErpId:         a.ErpID,
			AddressCode:   a.AddressCode,
			Name:          a.Name,
			Company:       a.Company,
			Address:       a.Address,
			PostalCode:    a.PostalCode,
			City:          a.City,
			CountryCode:   a.CountryCode,
			Phone:         a.Phone,
			Location:      a.Location,
			StreetType:    a.StreetType,
			Neighborhood:  a.Neighborhood,
			State:         a.State,
			CustomerErpId: a.CustomerErpID,
			OldErpId:      a.OldErpID,
			Migrated:      a.Migrated,
		})
	}
	return out
}
