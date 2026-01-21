package rest

import "time"

// DTOs kept in a single place to make transport shapes easy to find.
type customerInput struct {
	FrontendID                int64                  `json:"frontend_id"`
	ErpID                     string                 `json:"erp_id"`
	MarketID                  int64                  `json:"market_id"`
	MarketCustomizerID        int64                  `json:"market_customizer_id"`
	Level                     int32                  `json:"level"`
	ParentID                  int64                  `json:"parent_id"`
	DiscountProfile           string                 `json:"discount_profile"`
	IsActive                  bool                   `json:"is_active"`
	CanOrder                  bool                   `json:"can_order"`
	Username                  string                 `json:"username"`
	Password                  string                 `json:"password"`
	Email                     string                 `json:"email"`
	EmailCopy                 string                 `json:"email_copy"`
	CountryCode               string                 `json:"country_code"`
	Language                  string                 `json:"language"`
	ContactLanguage           string                 `json:"contact_language"`
	WebserviceKey             string                 `json:"webservice_key"`
	Name                      string                 `json:"name"`
	Company                   string                 `json:"company"`
	TaxID                     string                 `json:"tax_id"`
	Bank                      string                 `json:"bank"`
	BankAddress               string                 `json:"bank_address"`
	BankBranch                string                 `json:"bank_branch"`
	Website                   string                 `json:"website"`
	AddressLine1              string                 `json:"address_line1"`
	AddressLine2              string                 `json:"address_line2"`
	PostalCode                string                 `json:"postal_code"`
	City                      string                 `json:"city"`
	Phone                     string                 `json:"phone"`
	AccountManagerName        string                 `json:"account_manager_name"`
	AccountManagerPhone       string                 `json:"account_manager_phone"`
	AccountManagerEmail       string                 `json:"account_manager_email"`
	BirthDate                 *time.Time             `json:"birth_date,omitempty"`
	RegisteredAt              *time.Time             `json:"registered_at,omitempty"`
	LastLoginAt               *time.Time             `json:"last_login_at,omitempty"`
	FavoritesNotifications    bool                   `json:"favorites_notifications"`
	KeyCode                   int32                  `json:"key_code"`
	IsConfirmed               bool                   `json:"is_confirmed"`
	RecoveryTimestamp         *time.Time             `json:"recovery_timestamp,omitempty"`
	ReceivesNewsletters       bool                   `json:"receives_newsletters"`
	StandardTier              int32                  `json:"standard_tier"`
	OwnerID                   string                 `json:"owner_id"`
	StockPolicy               int32                  `json:"stock_policy"`
	Location                  string                 `json:"location"`
	StreetType                string                 `json:"street_type"`
	Neighborhood              string                 `json:"neighborhood"`
	State                     string                 `json:"state"`
	StateRegistration         string                 `json:"state_registration"`
	Country                   string                 `json:"country"`
	Comment                   string                 `json:"comment"`
	RegistrationCertificate   string                 `json:"registration_certificate"`
	LastLoginIP               string                 `json:"last_login_ip"`
	LastLoginCountryCode      string                 `json:"last_login_country_code"`
	BlockedBySuspiciousChange bool                   `json:"blocked_by_suspicious_change"`
	WarehouseCode             string                 `json:"warehouse_code"`
	OldErpID                  string                 `json:"old_erp_id"`
	CommercialMarketID        int64                  `json:"commercial_market_id"`
	Migrated                  bool                   `json:"migrated"`
	CommercialAreaID          string                 `json:"commercial_area_id"`
	IndustrialProduction      int32                  `json:"industrial_production"`
	DeliveryNote              string                 `json:"delivery_note"`
	NoDirectApprovals         bool                   `json:"no_direct_approvals"`
	IsCleaned                 bool                   `json:"is_cleaned"`
	Addresses                 []customerAddressInput `json:"addresses"`
}

type customerResponse struct {
	Code                      int64                     `json:"code"`
	FrontendID                int64                     `json:"frontend_id"`
	ErpID                     string                    `json:"erp_id"`
	MarketID                  int64                     `json:"market_id"`
	MarketCustomizerID        int64                     `json:"market_customizer_id"`
	Level                     int32                     `json:"level"`
	ParentID                  int64                     `json:"parent_id"`
	DiscountProfile           string                    `json:"discount_profile"`
	IsActive                  bool                      `json:"is_active"`
	CanOrder                  bool                      `json:"can_order"`
	Username                  string                    `json:"username"`
	Password                  string                    `json:"password"`
	Email                     string                    `json:"email"`
	EmailCopy                 string                    `json:"email_copy"`
	CountryCode               string                    `json:"country_code"`
	Language                  string                    `json:"language"`
	ContactLanguage           string                    `json:"contact_language"`
	WebserviceKey             string                    `json:"webservice_key"`
	Name                      string                    `json:"name"`
	Company                   string                    `json:"company"`
	TaxID                     string                    `json:"tax_id"`
	Bank                      string                    `json:"bank"`
	BankAddress               string                    `json:"bank_address"`
	BankBranch                string                    `json:"bank_branch"`
	Website                   string                    `json:"website"`
	AddressLine1              string                    `json:"address_line1"`
	AddressLine2              string                    `json:"address_line2"`
	PostalCode                string                    `json:"postal_code"`
	City                      string                    `json:"city"`
	Phone                     string                    `json:"phone"`
	AccountManagerName        string                    `json:"account_manager_name"`
	AccountManagerPhone       string                    `json:"account_manager_phone"`
	AccountManagerEmail       string                    `json:"account_manager_email"`
	BirthDate                 *time.Time                `json:"birth_date,omitempty"`
	RegisteredAt              *time.Time                `json:"registered_at,omitempty"`
	LastLoginAt               *time.Time                `json:"last_login_at,omitempty"`
	FavoritesNotifications    bool                      `json:"favorites_notifications"`
	KeyCode                   int32                     `json:"key_code"`
	IsConfirmed               bool                      `json:"is_confirmed"`
	RecoveryTimestamp         *time.Time                `json:"recovery_timestamp,omitempty"`
	ReceivesNewsletters       bool                      `json:"receives_newsletters"`
	StandardTier              int32                     `json:"standard_tier"`
	OwnerID                   string                    `json:"owner_id"`
	StockPolicy               int32                     `json:"stock_policy"`
	Location                  string                    `json:"location"`
	StreetType                string                    `json:"street_type"`
	Neighborhood              string                    `json:"neighborhood"`
	State                     string                    `json:"state"`
	StateRegistration         string                    `json:"state_registration"`
	Country                   string                    `json:"country"`
	Comment                   string                    `json:"comment"`
	RegistrationCertificate   string                    `json:"registration_certificate"`
	LastLoginIP               string                    `json:"last_login_ip"`
	LastLoginCountryCode      string                    `json:"last_login_country_code"`
	BlockedBySuspiciousChange bool                      `json:"blocked_by_suspicious_change"`
	WarehouseCode             string                    `json:"warehouse_code"`
	OldErpID                  string                    `json:"old_erp_id"`
	CommercialMarketID        int64                     `json:"commercial_market_id"`
	Migrated                  bool                      `json:"migrated"`
	CommercialAreaID          string                    `json:"commercial_area_id"`
	IndustrialProduction      int32                     `json:"industrial_production"`
	DeliveryNote              string                    `json:"delivery_note"`
	NoDirectApprovals         bool                      `json:"no_direct_approvals"`
	IsCleaned                 bool                      `json:"is_cleaned"`
	Addresses                 []customerAddressResponse `json:"addresses"`
	CreatedAt                 time.Time                 `json:"created_at,omitempty"`
	UpdatedAt                 time.Time                 `json:"updated_at,omitempty"`
}

type customerAddressInput struct {
	ID            int64  `json:"id"`
	ErpID         string `json:"erp_id"`
	AddressCode   string `json:"address_code"`
	Name          string `json:"name"`
	Company       string `json:"company"`
	Address       string `json:"address"`
	PostalCode    string `json:"postal_code"`
	City          string `json:"city"`
	CountryCode   string `json:"country_code"`
	Phone         string `json:"phone"`
	Location      string `json:"location"`
	StreetType    string `json:"street_type"`
	Neighborhood  string `json:"neighborhood"`
	State         string `json:"state"`
	CustomerErpID string `json:"customer_erp_id"`
	OldErpID      string `json:"old_erp_id"`
	Migrated      bool   `json:"migrated"`
}

type customerAddressResponse struct {
	ID            int64  `json:"id"`
	CustomerCode  int64  `json:"customer_code"`
	ErpID         string `json:"erp_id"`
	AddressCode   string `json:"address_code"`
	Name          string `json:"name"`
	Company       string `json:"company"`
	Address       string `json:"address"`
	PostalCode    string `json:"postal_code"`
	City          string `json:"city"`
	CountryCode   string `json:"country_code"`
	Phone         string `json:"phone"`
	Location      string `json:"location"`
	StreetType    string `json:"street_type"`
	Neighborhood  string `json:"neighborhood"`
	State         string `json:"state"`
	CustomerErpID string `json:"customer_erp_id"`
	OldErpID      string `json:"old_erp_id"`
	Migrated      bool   `json:"migrated"`
}

type pageResponse struct {
	NextPageToken string `json:"next_page_token"`
	TotalSize     int32  `json:"total_size"`
}

type listResponse struct {
	Customers []customerResponse `json:"customers"`
	Page      pageResponse       `json:"page"`
}
