package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	customerv1 "github.com/umbranian0/customer-mdm/api/gen/customer/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Gateway exposes the CustomerService gRPC API over a minimal REST/JSON surface.
// It is intentionally light-weight to avoid extra dependencies.
type Gateway struct {
	Client       customerv1.CustomerServiceClient
	CallTimeout  time.Duration
	MaxPageSize  int32
	DefaultLimit int32
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/customers") {
		http.NotFound(w, r)
		return
	}

	// Exact /customers route.
	if r.URL.Path == "/customers" {
		switch r.Method {
		case http.MethodPost:
			g.createCustomer(w, r)
		case http.MethodGet:
			g.listCustomers(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		}
		return
	}

	// Routes with an ID: /customers/{id}
	if strings.HasPrefix(r.URL.Path, "/customers/") {
		id := strings.TrimPrefix(r.URL.Path, "/customers/")
		if strings.Contains(id, "/") || id == "" {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			g.getCustomer(w, r, id)
		case http.MethodPut:
			g.updateCustomer(w, r, id)
		case http.MethodDelete:
			g.deleteCustomer(w, r, id)
		default:
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		}
		return
	}

	http.NotFound(w, r)
}

func (g *Gateway) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	timeout := g.CallTimeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return context.WithTimeout(ctx, timeout)
}

func (g *Gateway) createCustomer(w http.ResponseWriter, r *http.Request) {
	var in customerInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := g.withTimeout(r.Context())
	defer cancel()

	resp, err := g.Client.CreateCustomer(ctx, &customerv1.CreateCustomerRequest{
		Input:          toProtoInput(in),
		IdempotencyKey: r.Header.Get("Idempotency-Key"),
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toHTTPResponse(resp.Customer))
}

func (g *Gateway) getCustomer(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := g.withTimeout(r.Context())
	defer cancel()

	resp, err := g.Client.GetCustomer(ctx, &customerv1.GetCustomerRequest{Id: id})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toHTTPResponse(resp.Customer))
}

func (g *Gateway) updateCustomer(w http.ResponseWriter, r *http.Request, id string) {
	var in customerInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := g.withTimeout(r.Context())
	defer cancel()

	resp, err := g.Client.UpdateCustomer(ctx, &customerv1.UpdateCustomerRequest{
		Id:    id,
		Input: toProtoInput(in),
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toHTTPResponse(resp.Customer))
}

func (g *Gateway) deleteCustomer(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := g.withTimeout(r.Context())
	defer cancel()

	resp, err := g.Client.DeleteCustomer(ctx, &customerv1.DeleteCustomerRequest{Id: id})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"deleted": resp.GetDeleted()})
}

func (g *Gateway) listCustomers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	size := parsePageSize(q.Get("page_size"))
	if g.MaxPageSize > 0 && size > g.MaxPageSize {
		size = g.MaxPageSize
	}
	if size == 0 {
		size = g.DefaultLimit
		if size == 0 {
			size = 50
		}
	}

	ctx, cancel := g.withTimeout(r.Context())
	defer cancel()

	resp, err := g.Client.ListCustomers(ctx, &customerv1.ListCustomersRequest{
		Page: &customerv1.PageRequest{
			PageSize:  size,
			PageToken: q.Get("page_token"),
			Query:     q.Get("query"),
		},
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	out := listResponse{}
	if resp.GetPage() != nil {
		out.Page = pageResponse{
			NextPageToken: resp.GetPage().GetNextPageToken(),
			TotalSize:     resp.GetPage().GetTotalSize(),
		}
	}
	for _, c := range resp.GetCustomers() {
		out.Customers = append(out.Customers, toHTTPResponse(c))
	}
	writeJSON(w, http.StatusOK, out)
}

func parsePageSize(raw string) int32 {
	if raw == "" {
		return 0
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 0
	}
	return int32(n)
}

func toProtoInput(in customerInput) *customerv1.CustomerInput {
	return &customerv1.CustomerInput{
		Name:       in.Name,
		Email:      in.Email,
		TaxId:      in.TaxID,
		Phone:      in.Phone,
		Country:    in.Country,
		IsActive:   in.IsActive,
		Attributes: in.Attributes,
	}
}

func toHTTPResponse(c *customerv1.Customer) customerResponse {
	if c == nil {
		return customerResponse{}
	}
	return customerResponse{
		ID:         c.GetId(),
		Name:       c.GetName(),
		Email:      c.GetEmail(),
		TaxID:      c.GetTaxId(),
		Phone:      c.GetPhone(),
		Country:    c.GetCountry(),
		IsActive:   c.GetIsActive(),
		Attributes: c.GetAttributes(),
		CreatedAt:  toTime(c.GetCreatedAt()),
		UpdatedAt:  toTime(c.GetUpdatedAt()),
	}
}

func toTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func writeJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeError(w, grpcToHTTP(st.Code()), errors.New(st.Message()))
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	writeJSON(w, statusCode, map[string]string{"error": err.Error()})
}

func grpcToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.InvalidArgument, codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusBadGateway
	}
}
