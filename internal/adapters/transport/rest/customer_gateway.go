package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	customerv1 "github.com/umbranian0/customer-mdm/api/gen/api/proto/customer/v1"
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
	code, err := parseID(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx, cancel := g.withTimeout(r.Context())
	defer cancel()

	resp, err := g.Client.GetCustomer(ctx, &customerv1.GetCustomerRequest{Code: code})
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
	code, err := parseID(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := g.withTimeout(r.Context())
	defer cancel()

	resp, err := g.Client.UpdateCustomer(ctx, &customerv1.UpdateCustomerRequest{
		Code:  code,
		Input: toProtoInput(in),
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toHTTPResponse(resp.Customer))
}

func (g *Gateway) deleteCustomer(w http.ResponseWriter, r *http.Request, id string) {
	code, err := parseID(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx, cancel := g.withTimeout(r.Context())
	defer cancel()

	resp, err := g.Client.DeleteCustomer(ctx, &customerv1.DeleteCustomerRequest{Code: code})
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

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid customer code")
	}
	return id, nil
}

func toProtoAddressInputs(in []customerAddressInput) []*customerv1.CustomerAddressInput {
	if len(in) == 0 {
		return nil
	}
	out := make([]*customerv1.CustomerAddressInput, 0, len(in))
	for _, a := range in {
		out = append(out, &customerv1.CustomerAddressInput{
			Id:            a.ID,
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

func toProtoTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}

func toProtoInput(in customerInput) *customerv1.CustomerInput {
	return &customerv1.CustomerInput{
		FrontendId:                in.FrontendID,
		ErpId:                     in.ErpID,
		MarketId:                  in.MarketID,
		MarketCustomizerId:        in.MarketCustomizerID,
		Level:                     in.Level,
		ParentId:                  in.ParentID,
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
		TaxId:                     in.TaxID,
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
		BirthDate:                 toProtoTime(in.BirthDate),
		RegisteredAt:              toProtoTime(in.RegisteredAt),
		LastLoginAt:               toProtoTime(in.LastLoginAt),
		FavoritesNotifications:    in.FavoritesNotifications,
		KeyCode:                   in.KeyCode,
		IsConfirmed:               in.IsConfirmed,
		RecoveryTimestamp:         toProtoTime(in.RecoveryTimestamp),
		ReceivesNewsletters:       in.ReceivesNewsletters,
		StandardTier:              in.StandardTier,
		OwnerId:                   in.OwnerID,
		StockPolicy:               in.StockPolicy,
		Location:                  in.Location,
		StreetType:                in.StreetType,
		Neighborhood:              in.Neighborhood,
		State:                     in.State,
		StateRegistration:         in.StateRegistration,
		Country:                   in.Country,
		Comment:                   in.Comment,
		RegistrationCertificate:   in.RegistrationCertificate,
		LastLoginIp:               in.LastLoginIP,
		LastLoginCountryCode:      in.LastLoginCountryCode,
		BlockedBySuspiciousChange: in.BlockedBySuspiciousChange,
		WarehouseCode:             in.WarehouseCode,
		OldErpId:                  in.OldErpID,
		CommercialMarketId:        in.CommercialMarketID,
		Migrated:                  in.Migrated,
		CommercialAreaId:          in.CommercialAreaID,
		IndustrialProduction:      in.IndustrialProduction,
		DeliveryNote:              in.DeliveryNote,
		NoDirectApprovals:         in.NoDirectApprovals,
		IsCleaned:                 in.IsCleaned,
		Addresses:                 toProtoAddressInputs(in.Addresses),
	}
}

func toHTTPAddresses(in []*customerv1.CustomerAddress) []customerAddressResponse {
	if len(in) == 0 {
		return nil
	}
	out := make([]customerAddressResponse, 0, len(in))
	for _, a := range in {
		if a == nil {
			continue
		}
		out = append(out, customerAddressResponse{
			ID:            a.GetId(),
			CustomerCode:  a.GetCustomerCode(),
			ErpID:         a.GetErpId(),
			AddressCode:   a.GetAddressCode(),
			Name:          a.GetName(),
			Company:       a.GetCompany(),
			Address:       a.GetAddress(),
			PostalCode:    a.GetPostalCode(),
			City:          a.GetCity(),
			CountryCode:   a.GetCountryCode(),
			Phone:         a.GetPhone(),
			Location:      a.GetLocation(),
			StreetType:    a.GetStreetType(),
			Neighborhood:  a.GetNeighborhood(),
			State:         a.GetState(),
			CustomerErpID: a.GetCustomerErpId(),
			OldErpID:      a.GetOldErpId(),
			Migrated:      a.GetMigrated(),
		})
	}
	return out
}

func toHTTPResponse(c *customerv1.Customer) customerResponse {
	if c == nil {
		return customerResponse{}
	}
	return customerResponse{
		Code:                      c.GetCode(),
		FrontendID:                c.GetFrontendId(),
		ErpID:                     c.GetErpId(),
		MarketID:                  c.GetMarketId(),
		MarketCustomizerID:        c.GetMarketCustomizerId(),
		Level:                     c.GetLevel(),
		ParentID:                  c.GetParentId(),
		DiscountProfile:           c.GetDiscountProfile(),
		IsActive:                  c.GetIsActive(),
		CanOrder:                  c.GetCanOrder(),
		Username:                  c.GetUsername(),
		Password:                  c.GetPassword(),
		Email:                     c.GetEmail(),
		EmailCopy:                 c.GetEmailCopy(),
		CountryCode:               c.GetCountryCode(),
		Language:                  c.GetLanguage(),
		ContactLanguage:           c.GetContactLanguage(),
		WebserviceKey:             c.GetWebserviceKey(),
		Name:                      c.GetName(),
		Company:                   c.GetCompany(),
		TaxID:                     c.GetTaxId(),
		Bank:                      c.GetBank(),
		BankAddress:               c.GetBankAddress(),
		BankBranch:                c.GetBankBranch(),
		Website:                   c.GetWebsite(),
		AddressLine1:              c.GetAddressLine1(),
		AddressLine2:              c.GetAddressLine2(),
		PostalCode:                c.GetPostalCode(),
		City:                      c.GetCity(),
		Phone:                     c.GetPhone(),
		AccountManagerName:        c.GetAccountManagerName(),
		AccountManagerPhone:       c.GetAccountManagerPhone(),
		AccountManagerEmail:       c.GetAccountManagerEmail(),
		BirthDate:                 toTimePtr(c.GetBirthDate()),
		RegisteredAt:              toTimePtr(c.GetRegisteredAt()),
		LastLoginAt:               toTimePtr(c.GetLastLoginAt()),
		FavoritesNotifications:    c.GetFavoritesNotifications(),
		KeyCode:                   c.GetKeyCode(),
		IsConfirmed:               c.GetIsConfirmed(),
		RecoveryTimestamp:         toTimePtr(c.GetRecoveryTimestamp()),
		ReceivesNewsletters:       c.GetReceivesNewsletters(),
		StandardTier:              c.GetStandardTier(),
		OwnerID:                   c.GetOwnerId(),
		StockPolicy:               c.GetStockPolicy(),
		Location:                  c.GetLocation(),
		StreetType:                c.GetStreetType(),
		Neighborhood:              c.GetNeighborhood(),
		State:                     c.GetState(),
		StateRegistration:         c.GetStateRegistration(),
		Country:                   c.GetCountry(),
		Comment:                   c.GetComment(),
		RegistrationCertificate:   c.GetRegistrationCertificate(),
		LastLoginIP:               c.GetLastLoginIp(),
		LastLoginCountryCode:      c.GetLastLoginCountryCode(),
		BlockedBySuspiciousChange: c.GetBlockedBySuspiciousChange(),
		WarehouseCode:             c.GetWarehouseCode(),
		OldErpID:                  c.GetOldErpId(),
		CommercialMarketID:        c.GetCommercialMarketId(),
		Migrated:                  c.GetMigrated(),
		CommercialAreaID:          c.GetCommercialAreaId(),
		IndustrialProduction:      c.GetIndustrialProduction(),
		DeliveryNote:              c.GetDeliveryNote(),
		NoDirectApprovals:         c.GetNoDirectApprovals(),
		IsCleaned:                 c.GetIsCleaned(),
		Addresses:                 toHTTPAddresses(c.GetAddresses()),
		CreatedAt:                 toTime(c.GetCreatedAt()),
		UpdatedAt:                 toTime(c.GetUpdatedAt()),
	}
}

func toTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func toTimePtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
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
