CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS customers (
  code BIGSERIAL PRIMARY KEY,
  frontend_id BIGINT NOT NULL DEFAULT 0,
  erp_id TEXT NOT NULL DEFAULT '',
  market_id BIGINT NOT NULL DEFAULT 0,
  market_customizer_id BIGINT NOT NULL DEFAULT 0,
  level SMALLINT NOT NULL DEFAULT 0,
  parent_id BIGINT NOT NULL DEFAULT 0,
  discount_profile TEXT NOT NULL DEFAULT '',
  is_active BOOLEAN NOT NULL DEFAULT false,
  can_order BOOLEAN NOT NULL DEFAULT false,
  username TEXT NOT NULL DEFAULT '',
  password TEXT NOT NULL DEFAULT '',
  email TEXT,
  email_copy TEXT NOT NULL DEFAULT '',
  country_code TEXT NOT NULL DEFAULT '',
  language TEXT NOT NULL DEFAULT '',
  contact_language TEXT NOT NULL DEFAULT '',
  webservice_key TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL DEFAULT '',
  company TEXT NOT NULL DEFAULT '',
  tax_id TEXT NOT NULL DEFAULT '',
  bank TEXT NOT NULL DEFAULT '',
  bank_address TEXT NOT NULL DEFAULT '',
  bank_branch TEXT NOT NULL DEFAULT '',
  website TEXT NOT NULL DEFAULT '',
  address_line1 TEXT NOT NULL DEFAULT '',
  address_line2 TEXT NOT NULL DEFAULT '',
  postal_code TEXT NOT NULL DEFAULT '',
  city TEXT NOT NULL DEFAULT '',
  phone TEXT,
  account_manager_name TEXT NOT NULL DEFAULT '',
  account_manager_phone TEXT NOT NULL DEFAULT '',
  account_manager_email TEXT NOT NULL DEFAULT '',
  birth_date DATE,
  registered_at TIMESTAMPTZ,
  last_login_at TIMESTAMPTZ,
  favorites_notifications BOOLEAN NOT NULL DEFAULT true,
  key_code INT NOT NULL DEFAULT 0,
  is_confirmed BOOLEAN NOT NULL DEFAULT false,
  recovery_timestamp TIMESTAMPTZ,
  receives_newsletters BOOLEAN NOT NULL DEFAULT false,
  standard_tier INT NOT NULL DEFAULT 0,
  owner_id TEXT NOT NULL DEFAULT '',
  stock_policy INT NOT NULL DEFAULT 2,
  location TEXT NOT NULL DEFAULT '',
  street_type TEXT NOT NULL DEFAULT '',
  neighborhood TEXT NOT NULL DEFAULT '',
  state TEXT NOT NULL DEFAULT '',
  state_registration TEXT NOT NULL DEFAULT '',
  country TEXT,
  comment TEXT NOT NULL DEFAULT '',
  registration_certificate TEXT DEFAULT '',
  last_login_ip TEXT,
  last_login_country_code TEXT,
  blocked_by_suspicious_change BOOLEAN NOT NULL DEFAULT false,
  warehouse_code TEXT,
  old_erp_id TEXT,
  commercial_market_id BIGINT NOT NULL DEFAULT 0,
  migrated BOOLEAN NOT NULL DEFAULT true,
  commercial_area_id TEXT NOT NULL DEFAULT '',
  industrial_production INT NOT NULL DEFAULT 0,
  delivery_note TEXT DEFAULT '',
  no_direct_approvals BOOLEAN NOT NULL DEFAULT false,
  is_cleaned BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS customers_email_idx ON customers (email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS customers_phone_idx ON customers (phone) WHERE phone IS NOT NULL;
CREATE INDEX IF NOT EXISTS customers_erp_id_idx ON customers (erp_id);
CREATE INDEX IF NOT EXISTS customers_name_idx ON customers (name);
CREATE INDEX IF NOT EXISTS customers_registered_at_idx ON customers (registered_at);
CREATE INDEX IF NOT EXISTS customers_is_active_idx ON customers (is_active);
CREATE INDEX IF NOT EXISTS customers_market_id_idx ON customers (market_id);
CREATE INDEX IF NOT EXISTS customers_level_idx ON customers (level);
CREATE INDEX IF NOT EXISTS customers_parent_id_idx ON customers (parent_id);
CREATE INDEX IF NOT EXISTS customers_discount_profile_idx ON customers (discount_profile);

CREATE TABLE IF NOT EXISTS customer_addresses (
  id BIGSERIAL PRIMARY KEY,
  customer_code BIGINT NOT NULL REFERENCES customers(code) ON DELETE CASCADE,
  erp_id TEXT NOT NULL DEFAULT '',
  address_code TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL DEFAULT '',
  company TEXT NOT NULL DEFAULT '',
  address TEXT NOT NULL DEFAULT '',
  postal_code TEXT NOT NULL DEFAULT '',
  city TEXT NOT NULL DEFAULT '',
  country_code TEXT NOT NULL DEFAULT '',
  phone TEXT NOT NULL DEFAULT '',
  location TEXT NOT NULL DEFAULT '',
  street_type TEXT NOT NULL DEFAULT '',
  neighborhood TEXT NOT NULL DEFAULT '',
  state TEXT NOT NULL DEFAULT '',
  customer_erp_id TEXT,
  old_erp_id TEXT,
  migrated BOOLEAN NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS customer_addresses_customer_erp_idx
  ON customer_addresses (customer_code, erp_id);
CREATE INDEX IF NOT EXISTS customer_addresses_country_code_idx
  ON customer_addresses (country_code);
