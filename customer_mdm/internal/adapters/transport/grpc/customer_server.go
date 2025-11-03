package grpcadp

import (
	"context"
	"time"

	customerv1 "github.com/umbranian0/customer-mdm/api/gen/customer/v1"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/ports"
	"github.com/umbranian0/customer-mdm/internal/usecase"
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

func toProto(c *domain.Customer) *customerv1.Customer {
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

func fromInput(id string, in *customerv1.CustomerInput) *domain.Customer {
	now := time.Now().UTC()
	return &domain.Customer{
		ID:         id,
		Name:       in.Name,
		Email:      in.Email,
		TaxID:      in.TaxId,
		Phone:      in.Phone,
		Country:    in.Country,
		IsActive:   in.IsActive,
		Attributes: in.Attributes,
		UpdatedAt:  now,
	}
}

func (s *CustomerServer) CreateCustomer(ctx context.Context, req *customerv1.CreateCustomerRequest) (*customerv1.CreateCustomerResponse, error) {
	out, err := s.CreateUC.Run(ctx, usecase.CreateCustomerInput{
		Name: req.Input.Name, Email: req.Input.Email, TaxID: req.Input.TaxId,
		Phone: req.Input.Phone, Country: req.Input.Country, IsActive: req.Input.IsActive,
		Attributes: req.Input.Attributes, IdemKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	return &customerv1.CreateCustomerResponse{Customer: toProto(out)}, nil
}

func (s *CustomerServer) GetCustomer(ctx context.Context, req *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error) {
	out, err := s.GetUC.Run(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &customerv1.GetCustomerResponse{Customer: toProto(out)}, nil
}

func (s *CustomerServer) UpdateCustomer(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error) {
	out, err := s.UpdateUC.Run(ctx, usecase.UpdateCustomerInput{
		ID:   req.Id,
		Name: req.Input.Name, Email: req.Input.Email, TaxID: req.Input.TaxId,
		Phone: req.Input.Phone, Country: req.Input.Country, IsActive: req.Input.IsActive,
		Attributes: req.Input.Attributes,
	})
	if err != nil {
		return nil, err
	}
	return &customerv1.UpdateCustomerResponse{Customer: toProto(out)}, nil
}

func (s *CustomerServer) DeleteCustomer(ctx context.Context, req *customerv1.DeleteCustomerRequest) (*customerv1.DeleteCustomerResponse, error) {
	err := s.DeleteUC.Run(ctx, req.Id)
	if err != nil {
		return nil, err
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

var _ ports.EventPublisher = (*dummyPublisher)(nil)

type dummyPublisher struct{}

func (d *dummyPublisher) Publish(ctx context.Context, ev ports.Event) error { return nil }
