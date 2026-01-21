package grpcadp

import (
	"context"
	"errors"
	"time"

	customerv1 "github.com/umbranian0/customer-mdm/api/gen/api/proto/customer/v1"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/ports"
	"github.com/umbranian0/customer-mdm/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CustomerServer struct {
	customerv1.UnimplementedCustomerServiceServer
	CreateUC *usecase.CreateCustomer
	GetUC    *usecase.GetCustomer
	UpdateUC *usecase.UpdateCustomer
	DeleteUC *usecase.DeleteCustomer
	ListUC   *usecase.ListCustomers
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

func fromTimestamp(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func toProtoAddress(a domain.CustomerAddress) *customerv1.CustomerAddress {
	return &customerv1.CustomerAddress{
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
	}
}

func toProtoAddresses(addresses []domain.CustomerAddress) []*customerv1.CustomerAddress {
	if len(addresses) == 0 {
		return nil
	}
	out := make([]*customerv1.CustomerAddress, 0, len(addresses))
	for _, a := range addresses {
		out = append(out, toProtoAddress(a))
	}
	return out
}

func toProto(c *domain.Customer) *customerv1.Customer {
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

func fromAddressInputs(code int64, in []*customerv1.CustomerAddressInput) []domain.CustomerAddress {
	if len(in) == 0 {
		return nil
	}
	out := make([]domain.CustomerAddress, 0, len(in))
	for _, a := range in {
		if a == nil {
			continue
		}
		out = append(out, domain.CustomerAddress{
			ID:            a.Id,
			CustomerCode:  code,
			ErpID:         a.ErpId,
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
			CustomerErpID: a.CustomerErpId,
			OldErpID:      a.OldErpId,
			Migrated:      a.Migrated,
		})
	}
	return out
}

func fromInput(code int64, in *customerv1.CustomerInput) *domain.Customer {
	if in == nil {
		return &domain.Customer{Code: code}
	}
	return &domain.Customer{
		Code:                      code,
		FrontendID:                in.FrontendId,
		ErpID:                     in.ErpId,
		MarketID:                  in.MarketId,
		MarketCustomizerID:        in.MarketCustomizerId,
		Level:                     in.Level,
		ParentID:                  in.ParentId,
		DiscountProfile:           in.DiscountProfile,
		IsActive:                  in.IsActive,
		CanOrder:                  in.CanOrder,
		Username:                  in.Username,
		Password:                  in.Password,
		Email:                     in.Email,
		EmailCopy:                 in.EmailCopy,
		CountryCode:               in.CountryCode,
		Language:                  in.Language,
		ContactLanguage:           in.ContactLanguage,
		WebserviceKey:             in.WebserviceKey,
		Name:                      in.Name,
		Company:                   in.Company,
		TaxID:                     in.TaxId,
		Bank:                      in.Bank,
		BankAddress:               in.BankAddress,
		BankBranch:                in.BankBranch,
		Website:                   in.Website,
		AddressLine1:              in.AddressLine1,
		AddressLine2:              in.AddressLine2,
		PostalCode:                in.PostalCode,
		City:                      in.City,
		Phone:                     in.Phone,
		AccountManagerName:        in.AccountManagerName,
		AccountManagerPhone:       in.AccountManagerPhone,
		AccountManagerEmail:       in.AccountManagerEmail,
		BirthDate:                 fromTimestamp(in.BirthDate),
		RegisteredAt:              fromTimestamp(in.RegisteredAt),
		LastLoginAt:               fromTimestamp(in.LastLoginAt),
		FavoritesNotifications:    in.FavoritesNotifications,
		KeyCode:                   in.KeyCode,
		IsConfirmed:               in.IsConfirmed,
		RecoveryTimestamp:         fromTimestamp(in.RecoveryTimestamp),
		ReceivesNewsletters:       in.ReceivesNewsletters,
		StandardTier:              in.StandardTier,
		OwnerID:                   in.OwnerId,
		StockPolicy:               in.StockPolicy,
		Location:                  in.Location,
		StreetType:                in.StreetType,
		Neighborhood:              in.Neighborhood,
		State:                     in.State,
		StateRegistration:         in.StateRegistration,
		Country:                   in.Country,
		Comment:                   in.Comment,
		RegistrationCertificate:   in.RegistrationCertificate,
		LastLoginIP:               in.LastLoginIp,
		LastLoginCountryCode:      in.LastLoginCountryCode,
		BlockedBySuspiciousChange: in.BlockedBySuspiciousChange,
		WarehouseCode:             in.WarehouseCode,
		OldErpID:                  in.OldErpId,
		CommercialMarketID:        in.CommercialMarketId,
		Migrated:                  in.Migrated,
		CommercialAreaID:          in.CommercialAreaId,
		IndustrialProduction:      in.IndustrialProduction,
		DeliveryNote:              in.DeliveryNote,
		NoDirectApprovals:         in.NoDirectApprovals,
		IsCleaned:                 in.IsCleaned,
		Addresses:                 fromAddressInputs(code, in.Addresses),
	}
}

func (s *CustomerServer) CreateCustomer(ctx context.Context, req *customerv1.CreateCustomerRequest) (*customerv1.CreateCustomerResponse, error) {
	out, err := s.CreateUC.Run(ctx, usecase.CreateCustomerInput{
		Customer: fromInput(0, req.Input),
		IdemKey:  req.IdempotencyKey,
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return &customerv1.CreateCustomerResponse{Customer: toProto(out)}, nil
}

func (s *CustomerServer) GetCustomer(ctx context.Context, req *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error) {
	out, err := s.GetUC.Run(ctx, req.Code)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &customerv1.GetCustomerResponse{Customer: toProto(out)}, nil
}

func (s *CustomerServer) UpdateCustomer(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error) {
	out, err := s.UpdateUC.Run(ctx, usecase.UpdateCustomerInput{
		Code:     req.Code,
		Customer: fromInput(req.Code, req.Input),
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return &customerv1.UpdateCustomerResponse{Customer: toProto(out)}, nil
}

func (s *CustomerServer) DeleteCustomer(ctx context.Context, req *customerv1.DeleteCustomerRequest) (*customerv1.DeleteCustomerResponse, error) {
	err := s.DeleteUC.Run(ctx, req.Code)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &customerv1.DeleteCustomerResponse{Deleted: true}, nil
}

func (s *CustomerServer) ListCustomers(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error) {
	items, next, total, err := s.ListUC.Run(ctx, usecase.ListCustomersInput{
		PageSize: int(req.Page.PageSize), PageToken: req.Page.PageToken, Query: req.Page.Query,
	})
	if err != nil {
		return nil, err
	}
	res := &customerv1.ListCustomersResponse{Page: &customerv1.PageResponse{
		NextPageToken: next, TotalSize: int32(total),
	}}
	for _, it := range items {
		res.Customers = append(res.Customers, toProto(it))
	}
	return res, nil
}

// DummyPublisher is a no-op EventPublisher used in wiring for local runs.
var _ ports.EventPublisher = (*DummyPublisher)(nil)

type DummyPublisher struct{}

func (d *DummyPublisher) Publish(ctx context.Context, ev ports.Event) error { return nil }

func toStatusError(err error) error {
	if err == nil {
		return nil
	}
	var conflict *domain.ConflictError
	if errors.As(err, &conflict) {
		return status.Error(codes.AlreadyExists, conflict.Error())
	}
	switch {
	case errors.Is(err, domain.ErrCustomerNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrCustomerConflict):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrCustomerBadInput):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
