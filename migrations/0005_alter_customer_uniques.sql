ALTER TABLE customers
  ALTER COLUMN email DROP NOT NULL,
  ALTER COLUMN email DROP DEFAULT,
  ALTER COLUMN phone DROP NOT NULL,
  ALTER COLUMN phone DROP DEFAULT;

UPDATE customers SET email = NULL WHERE email = '';
UPDATE customers SET phone = NULL WHERE phone = '';

DROP INDEX IF EXISTS customers_username_idx;
DROP INDEX IF EXISTS customers_email_idx;

CREATE UNIQUE INDEX IF NOT EXISTS customers_email_idx ON customers (email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS customers_phone_idx ON customers (phone) WHERE phone IS NOT NULL;
