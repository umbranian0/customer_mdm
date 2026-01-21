package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umbranian0/customer-mdm/internal/domain"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

type CustomerRepository struct {
	Pool *pgxpool.Pool
}

type OutboxWriter struct {
	Pool *pgxpool.Pool
}

func (r *CustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	if c == nil {
		return nil
	}
	now := time.Now().UTC()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = now
	}

	var (
		code int64
		err  error
	)
	if c.Code == 0 {
		const q = `INSERT INTO customers (
  frontend_id, erp_id, market_id, market_customizer_id, level, parent_id, discount_profile,
  is_active, can_order, username, password, email, email_copy, country_code, language,
  contact_language, webservice_key, name, company, tax_id, bank, bank_address, bank_branch,
  website, address_line1, address_line2, postal_code, city, phone, account_manager_name,
  account_manager_phone, account_manager_email, birth_date, registered_at, last_login_at,
  favorites_notifications, key_code, is_confirmed, recovery_timestamp, receives_newsletters,
  standard_tier, owner_id, stock_policy, location, street_type, neighborhood, state,
  state_registration, country, comment, registration_certificate, last_login_ip,
  last_login_country_code, blocked_by_suspicious_change, warehouse_code, old_erp_id,
  commercial_market_id, migrated, commercial_area_id, industrial_production, delivery_note,
  no_direct_approvals, is_cleaned, created_at, updated_at
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,
  $8,$9,$10,$11,$12,$13,$14,$15,
  $16,$17,$18,$19,$20,$21,$22,$23,
  $24,$25,$26,$27,$28,$29,$30,$31,
  $32,$33,$34,$35,$36,$37,$38,$39,$40,
  $41,$42,$43,$44,$45,$46,$47,$48,
  $49,$50,$51,$52,$53,$54,$55,$56,$57,
  $58,$59,$60,$61,$62,$63,$64,$65
) RETURNING code`
		err = r.Pool.QueryRow(ctx, q,
			c.FrontendID, c.ErpID, c.MarketID, c.MarketCustomizerID, c.Level, c.ParentID, c.DiscountProfile,
			c.IsActive, c.CanOrder, c.Username, c.Password, nullIfEmpty(c.Email), c.EmailCopy, c.CountryCode, c.Language,
			c.ContactLanguage, c.WebserviceKey, c.Name, c.Company, c.TaxID, c.Bank, c.BankAddress, c.BankBranch,
			c.Website, c.AddressLine1, c.AddressLine2, c.PostalCode, c.City, nullIfEmpty(c.Phone), c.AccountManagerName,
			c.AccountManagerPhone, c.AccountManagerEmail, c.BirthDate, c.RegisteredAt, c.LastLoginAt,
			c.FavoritesNotifications, c.KeyCode, c.IsConfirmed, c.RecoveryTimestamp, c.ReceivesNewsletters,
			c.StandardTier, c.OwnerID, c.StockPolicy, c.Location, c.StreetType, c.Neighborhood, c.State,
			c.StateRegistration, c.Country, c.Comment, c.RegistrationCertificate, c.LastLoginIP,
			c.LastLoginCountryCode, c.BlockedBySuspiciousChange, c.WarehouseCode, c.OldErpID,
			c.CommercialMarketID, c.Migrated, c.CommercialAreaID, c.IndustrialProduction, c.DeliveryNote,
			c.NoDirectApprovals, c.IsCleaned, c.CreatedAt, c.UpdatedAt,
		).Scan(&code)
	} else {
		const q = `INSERT INTO customers (
  code, frontend_id, erp_id, market_id, market_customizer_id, level, parent_id, discount_profile,
  is_active, can_order, username, password, email, email_copy, country_code, language,
  contact_language, webservice_key, name, company, tax_id, bank, bank_address, bank_branch,
  website, address_line1, address_line2, postal_code, city, phone, account_manager_name,
  account_manager_phone, account_manager_email, birth_date, registered_at, last_login_at,
  favorites_notifications, key_code, is_confirmed, recovery_timestamp, receives_newsletters,
  standard_tier, owner_id, stock_policy, location, street_type, neighborhood, state,
  state_registration, country, comment, registration_certificate, last_login_ip,
  last_login_country_code, blocked_by_suspicious_change, warehouse_code, old_erp_id,
  commercial_market_id, migrated, commercial_area_id, industrial_production, delivery_note,
  no_direct_approvals, is_cleaned, created_at, updated_at
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,
  $9,$10,$11,$12,$13,$14,$15,$16,
  $17,$18,$19,$20,$21,$22,$23,$24,
  $25,$26,$27,$28,$29,$30,$31,$32,
  $33,$34,$35,$36,$37,$38,$39,$40,$41,
  $42,$43,$44,$45,$46,$47,$48,$49,
  $50,$51,$52,$53,$54,$55,$56,$57,$58,
  $59,$60,$61,$62,$63,$64,$65,$66
) RETURNING code`
		err = r.Pool.QueryRow(ctx, q,
			c.Code, c.FrontendID, c.ErpID, c.MarketID, c.MarketCustomizerID, c.Level, c.ParentID, c.DiscountProfile,
			c.IsActive, c.CanOrder, c.Username, c.Password, nullIfEmpty(c.Email), c.EmailCopy, c.CountryCode, c.Language,
			c.ContactLanguage, c.WebserviceKey, c.Name, c.Company, c.TaxID, c.Bank, c.BankAddress, c.BankBranch,
			c.Website, c.AddressLine1, c.AddressLine2, c.PostalCode, c.City, nullIfEmpty(c.Phone), c.AccountManagerName,
			c.AccountManagerPhone, c.AccountManagerEmail, c.BirthDate, c.RegisteredAt, c.LastLoginAt,
			c.FavoritesNotifications, c.KeyCode, c.IsConfirmed, c.RecoveryTimestamp, c.ReceivesNewsletters,
			c.StandardTier, c.OwnerID, c.StockPolicy, c.Location, c.StreetType, c.Neighborhood, c.State,
			c.StateRegistration, c.Country, c.Comment, c.RegistrationCertificate, c.LastLoginIP,
			c.LastLoginCountryCode, c.BlockedBySuspiciousChange, c.WarehouseCode, c.OldErpID,
			c.CommercialMarketID, c.Migrated, c.CommercialAreaID, c.IndustrialProduction, c.DeliveryNote,
			c.NoDirectApprovals, c.IsCleaned, c.CreatedAt, c.UpdatedAt,
		).Scan(&code)
	}
	if err != nil {
		return err
	}
	c.Code = code
	return r.replaceAddresses(ctx, c.Code, c.Addresses)
}

func (r *CustomerRepository) Get(ctx context.Context, code int64) (*domain.Customer, error) {
	const q = `SELECT code, frontend_id, erp_id, market_id, market_customizer_id, level, parent_id, discount_profile,
  is_active, can_order, username, password, email, email_copy, country_code, language, contact_language,
  webservice_key, name, company, tax_id, bank, bank_address, bank_branch, website, address_line1,
  address_line2, postal_code, city, phone, account_manager_name, account_manager_phone, account_manager_email,
  birth_date, registered_at, last_login_at, favorites_notifications, key_code, is_confirmed, recovery_timestamp,
  receives_newsletters, standard_tier, owner_id, stock_policy, location, street_type, neighborhood, state,
  state_registration, country, comment, registration_certificate, last_login_ip, last_login_country_code,
  blocked_by_suspicious_change, warehouse_code, old_erp_id, commercial_market_id, migrated, commercial_area_id,
  industrial_production, delivery_note, no_direct_approvals, is_cleaned, created_at, updated_at
FROM customers WHERE code=$1`
	row := r.Pool.QueryRow(ctx, q, code)
	return r.scanCustomer(ctx, row)
}

func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	const q = `SELECT code, frontend_id, erp_id, market_id, market_customizer_id, level, parent_id, discount_profile,
  is_active, can_order, username, password, email, email_copy, country_code, language, contact_language,
  webservice_key, name, company, tax_id, bank, bank_address, bank_branch, website, address_line1,
  address_line2, postal_code, city, phone, account_manager_name, account_manager_phone, account_manager_email,
  birth_date, registered_at, last_login_at, favorites_notifications, key_code, is_confirmed, recovery_timestamp,
  receives_newsletters, standard_tier, owner_id, stock_policy, location, street_type, neighborhood, state,
  state_registration, country, comment, registration_certificate, last_login_ip, last_login_country_code,
  blocked_by_suspicious_change, warehouse_code, old_erp_id, commercial_market_id, migrated, commercial_area_id,
  industrial_production, delivery_note, no_direct_approvals, is_cleaned, created_at, updated_at
FROM customers WHERE email=$1`
	row := r.Pool.QueryRow(ctx, q, email)
	return r.scanCustomer(ctx, row)
}

func (r *CustomerRepository) GetByPhone(ctx context.Context, phone string) (*domain.Customer, error) {
	const q = `SELECT code, frontend_id, erp_id, market_id, market_customizer_id, level, parent_id, discount_profile,
  is_active, can_order, username, password, email, email_copy, country_code, language, contact_language,
  webservice_key, name, company, tax_id, bank, bank_address, bank_branch, website, address_line1,
  address_line2, postal_code, city, phone, account_manager_name, account_manager_phone, account_manager_email,
  birth_date, registered_at, last_login_at, favorites_notifications, key_code, is_confirmed, recovery_timestamp,
  receives_newsletters, standard_tier, owner_id, stock_policy, location, street_type, neighborhood, state,
  state_registration, country, comment, registration_certificate, last_login_ip, last_login_country_code,
  blocked_by_suspicious_change, warehouse_code, old_erp_id, commercial_market_id, migrated, commercial_area_id,
  industrial_production, delivery_note, no_direct_approvals, is_cleaned, created_at, updated_at
FROM customers WHERE phone=$1`
	row := r.Pool.QueryRow(ctx, q, phone)
	return r.scanCustomer(ctx, row)
}

type customerRow interface {
	Scan(dest ...any) error
}

func (r *CustomerRepository) scanCustomer(ctx context.Context, row customerRow) (*domain.Customer, error) {
	var (
		c                 domain.Customer
		birthDate         pgtype.Date
		registeredAt      pgtype.Timestamptz
		lastLoginAt       pgtype.Timestamptz
		recoveryTimestamp pgtype.Timestamptz
		createdAt         time.Time
		updatedAt         time.Time
	)
	if err := row.Scan(
		&c.Code, &c.FrontendID, &c.ErpID, &c.MarketID, &c.MarketCustomizerID, &c.Level, &c.ParentID, &c.DiscountProfile,
		&c.IsActive, &c.CanOrder, &c.Username, &c.Password, &c.Email, &c.EmailCopy, &c.CountryCode, &c.Language, &c.ContactLanguage,
		&c.WebserviceKey, &c.Name, &c.Company, &c.TaxID, &c.Bank, &c.BankAddress, &c.BankBranch, &c.Website, &c.AddressLine1,
		&c.AddressLine2, &c.PostalCode, &c.City, &c.Phone, &c.AccountManagerName, &c.AccountManagerPhone, &c.AccountManagerEmail,
		&birthDate, &registeredAt, &lastLoginAt, &c.FavoritesNotifications, &c.KeyCode, &c.IsConfirmed, &recoveryTimestamp,
		&c.ReceivesNewsletters, &c.StandardTier, &c.OwnerID, &c.StockPolicy, &c.Location, &c.StreetType, &c.Neighborhood, &c.State,
		&c.StateRegistration, &c.Country, &c.Comment, &c.RegistrationCertificate, &c.LastLoginIP, &c.LastLoginCountryCode,
		&c.BlockedBySuspiciousChange, &c.WarehouseCode, &c.OldErpID, &c.CommercialMarketID, &c.Migrated, &c.CommercialAreaID,
		&c.IndustrialProduction, &c.DeliveryNote, &c.NoDirectApprovals, &c.IsCleaned, &createdAt, &updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	c.CreatedAt = createdAt
	c.UpdatedAt = updatedAt
	c.BirthDate = timePtrFromDate(birthDate)
	c.RegisteredAt = timePtrFromTimestamptz(registeredAt)
	c.LastLoginAt = timePtrFromTimestamptz(lastLoginAt)
	c.RecoveryTimestamp = timePtrFromTimestamptz(recoveryTimestamp)

	addresses, err := r.listAddresses(ctx, c.Code)
	if err != nil {
		return nil, err
	}
	c.Addresses = addresses
	return &c, nil
}

func (r *CustomerRepository) Update(ctx context.Context, c *domain.Customer) error {
	if c == nil {
		return nil
	}
	c.UpdatedAt = time.Now().UTC()
	const q = `UPDATE customers SET
  frontend_id=$2, erp_id=$3, market_id=$4, market_customizer_id=$5, level=$6, parent_id=$7, discount_profile=$8,
  is_active=$9, can_order=$10, username=$11, password=$12, email=$13, email_copy=$14, country_code=$15, language=$16,
  contact_language=$17, webservice_key=$18, name=$19, company=$20, tax_id=$21, bank=$22, bank_address=$23, bank_branch=$24,
  website=$25, address_line1=$26, address_line2=$27, postal_code=$28, city=$29, phone=$30, account_manager_name=$31,
  account_manager_phone=$32, account_manager_email=$33, birth_date=$34, registered_at=$35, last_login_at=$36,
  favorites_notifications=$37, key_code=$38, is_confirmed=$39, recovery_timestamp=$40, receives_newsletters=$41,
  standard_tier=$42, owner_id=$43, stock_policy=$44, location=$45, street_type=$46, neighborhood=$47, state=$48,
  state_registration=$49, country=$50, comment=$51, registration_certificate=$52, last_login_ip=$53,
  last_login_country_code=$54, blocked_by_suspicious_change=$55, warehouse_code=$56, old_erp_id=$57,
  commercial_market_id=$58, migrated=$59, commercial_area_id=$60, industrial_production=$61, delivery_note=$62,
  no_direct_approvals=$63, is_cleaned=$64, updated_at=$65
WHERE code=$1`
	cmd, err := r.Pool.Exec(ctx, q,
		c.Code, c.FrontendID, c.ErpID, c.MarketID, c.MarketCustomizerID, c.Level, c.ParentID, c.DiscountProfile,
		c.IsActive, c.CanOrder, c.Username, c.Password, nullIfEmpty(c.Email), c.EmailCopy, c.CountryCode, c.Language, c.ContactLanguage,
		c.WebserviceKey, c.Name, c.Company, c.TaxID, c.Bank, c.BankAddress, c.BankBranch, c.Website, c.AddressLine1,
		c.AddressLine2, c.PostalCode, c.City, nullIfEmpty(c.Phone), c.AccountManagerName, c.AccountManagerPhone, c.AccountManagerEmail,
		c.BirthDate, c.RegisteredAt, c.LastLoginAt, c.FavoritesNotifications, c.KeyCode, c.IsConfirmed, c.RecoveryTimestamp,
		c.ReceivesNewsletters, c.StandardTier, c.OwnerID, c.StockPolicy, c.Location, c.StreetType, c.Neighborhood, c.State,
		c.StateRegistration, c.Country, c.Comment, c.RegistrationCertificate, c.LastLoginIP, c.LastLoginCountryCode,
		c.BlockedBySuspiciousChange, c.WarehouseCode, c.OldErpID, c.CommercialMarketID, c.Migrated, c.CommercialAreaID,
		c.IndustrialProduction, c.DeliveryNote, c.NoDirectApprovals, c.IsCleaned, c.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return r.replaceAddresses(ctx, c.Code, c.Addresses)
}

func (r *CustomerRepository) Delete(ctx context.Context, code int64) error {
	const q = `DELETE FROM customers WHERE code=$1`
	cmd, err := r.Pool.Exec(ctx, q, code)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *CustomerRepository) List(ctx context.Context, pageSize int, pageToken, query string) (items []*domain.Customer, next string, total int, err error) {
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 500 {
		pageSize = 500
	}

	where := ""
	args := []any{pageSize}
	if q := strings.TrimSpace(query); q != "" {
		where = "WHERE name ILIKE '%' || $2 || '%' OR email ILIKE '%' || $2 || '%' OR erp_id ILIKE '%' || $2 || '%'"
		args = append([]any{pageSize, q}, args[1:]...)
	}
	sql := `SELECT code, frontend_id, erp_id, market_id, market_customizer_id, level, parent_id, discount_profile,
  is_active, can_order, username, password, email, email_copy, country_code, language, contact_language,
  webservice_key, name, company, tax_id, bank, bank_address, bank_branch, website, address_line1,
  address_line2, postal_code, city, phone, account_manager_name, account_manager_phone, account_manager_email,
  birth_date, registered_at, last_login_at, favorites_notifications, key_code, is_confirmed, recovery_timestamp,
  receives_newsletters, standard_tier, owner_id, stock_policy, location, street_type, neighborhood, state,
  state_registration, country, comment, registration_certificate, last_login_ip, last_login_country_code,
  blocked_by_suspicious_change, warehouse_code, old_erp_id, commercial_market_id, migrated, commercial_area_id,
  industrial_production, delivery_note, no_direct_approvals, is_cleaned, created_at, updated_at
FROM customers ` + where + `
ORDER BY created_at DESC
LIMIT $1`

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			c                 domain.Customer
			birthDate         pgtype.Date
			registeredAt      pgtype.Timestamptz
			lastLoginAt       pgtype.Timestamptz
			recoveryTimestamp pgtype.Timestamptz
			createdAt         time.Time
			updatedAt         time.Time
		)
		if err := rows.Scan(
			&c.Code, &c.FrontendID, &c.ErpID, &c.MarketID, &c.MarketCustomizerID, &c.Level, &c.ParentID, &c.DiscountProfile,
			&c.IsActive, &c.CanOrder, &c.Username, &c.Password, &c.Email, &c.EmailCopy, &c.CountryCode, &c.Language, &c.ContactLanguage,
			&c.WebserviceKey, &c.Name, &c.Company, &c.TaxID, &c.Bank, &c.BankAddress, &c.BankBranch, &c.Website, &c.AddressLine1,
			&c.AddressLine2, &c.PostalCode, &c.City, &c.Phone, &c.AccountManagerName, &c.AccountManagerPhone, &c.AccountManagerEmail,
			&birthDate, &registeredAt, &lastLoginAt, &c.FavoritesNotifications, &c.KeyCode, &c.IsConfirmed, &recoveryTimestamp,
			&c.ReceivesNewsletters, &c.StandardTier, &c.OwnerID, &c.StockPolicy, &c.Location, &c.StreetType, &c.Neighborhood, &c.State,
			&c.StateRegistration, &c.Country, &c.Comment, &c.RegistrationCertificate, &c.LastLoginIP, &c.LastLoginCountryCode,
			&c.BlockedBySuspiciousChange, &c.WarehouseCode, &c.OldErpID, &c.CommercialMarketID, &c.Migrated, &c.CommercialAreaID,
			&c.IndustrialProduction, &c.DeliveryNote, &c.NoDirectApprovals, &c.IsCleaned, &createdAt, &updatedAt,
		); err != nil {
			return nil, "", 0, err
		}
		c.CreatedAt = createdAt
		c.UpdatedAt = updatedAt
		c.BirthDate = timePtrFromDate(birthDate)
		c.RegisteredAt = timePtrFromTimestamptz(registeredAt)
		c.LastLoginAt = timePtrFromTimestamptz(lastLoginAt)
		c.RecoveryTimestamp = timePtrFromTimestamptz(recoveryTimestamp)
		items = append(items, &c)
	}
	_ = r.Pool.QueryRow(ctx, "SELECT COUNT(1) FROM customers").Scan(&total)
	next = ""
	return
}

func (o *OutboxWriter) Write(tx ports.Tx, topic string, key, value []byte, headers map[string]string) error {
	const q = `INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, payload, headers)
               VALUES ($1,$2,$3,$4,$5)`
	aggregateType := "customer"
	aggregateID := ""
	if headers != nil {
		aggregateID = strings.TrimSpace(headers["ID"])
	}
	if aggregateID == "" && len(key) > 0 {
		aggregateID = strings.TrimSpace(string(key))
	}
	if aggregateID == "" {
		aggregateID = uuid.NewString()
	}
	eventType := "generic"
	if v, ok := headers["event_type"]; ok && v != "" {
		eventType = v
	}
	if headers == nil {
		headers = map[string]string{}
	}
	_, err := o.Pool.Exec(tx.Context(), q, aggregateType, aggregateID, eventType, value, headers)
	return err
}

func (r *CustomerRepository) replaceAddresses(ctx context.Context, code int64, addresses []domain.CustomerAddress) error {
	if code == 0 {
		return nil
	}
	_, err := r.Pool.Exec(ctx, `DELETE FROM customer_addresses WHERE customer_code=$1`, code)
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return nil
	}
	for i := range addresses {
		addresses[i].CustomerCode = code
	}
	return r.insertAddresses(ctx, addresses)
}

func (r *CustomerRepository) insertAddresses(ctx context.Context, addresses []domain.CustomerAddress) error {
	const q = `INSERT INTO customer_addresses (
  customer_code, erp_id, address_code, name, company, address, postal_code, city, country_code, phone,
  location, street_type, neighborhood, state, customer_erp_id, old_erp_id, migrated
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	for _, a := range addresses {
		if _, err := r.Pool.Exec(ctx, q,
			a.CustomerCode, a.ErpID, a.AddressCode, a.Name, a.Company, a.Address, a.PostalCode, a.City, a.CountryCode, a.Phone,
			a.Location, a.StreetType, a.Neighborhood, a.State, a.CustomerErpID, a.OldErpID, a.Migrated,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *CustomerRepository) listAddresses(ctx context.Context, code int64) ([]domain.CustomerAddress, error) {
	const q = `SELECT id, customer_code, erp_id, address_code, name, company, address, postal_code, city, country_code,
  phone, location, street_type, neighborhood, state, customer_erp_id, old_erp_id, migrated
FROM customer_addresses WHERE customer_code=$1 ORDER BY id`
	rows, err := r.Pool.Query(ctx, q, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.CustomerAddress
	for rows.Next() {
		var a domain.CustomerAddress
		if err := rows.Scan(
			&a.ID, &a.CustomerCode, &a.ErpID, &a.AddressCode, &a.Name, &a.Company, &a.Address, &a.PostalCode, &a.City, &a.CountryCode,
			&a.Phone, &a.Location, &a.StreetType, &a.Neighborhood, &a.State, &a.CustomerErpID, &a.OldErpID, &a.Migrated,
		); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func timePtrFromDate(v pgtype.Date) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

func timePtrFromTimestamptz(v pgtype.Timestamptz) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
