WITH ranked AS (
  SELECT code, email,
         ROW_NUMBER() OVER (PARTITION BY email ORDER BY code) AS rn
  FROM customers
  WHERE email IS NOT NULL
),
email_dupes AS (
  SELECT code FROM ranked WHERE rn > 1
)
UPDATE customers
SET email = NULL
WHERE code IN (SELECT code FROM email_dupes);

WITH ranked AS (
  SELECT code, phone,
         ROW_NUMBER() OVER (PARTITION BY phone ORDER BY code) AS rn
  FROM customers
  WHERE phone IS NOT NULL
),
phone_dupes AS (
  SELECT code FROM ranked WHERE rn > 1
)
UPDATE customers
SET phone = NULL
WHERE code IN (SELECT code FROM phone_dupes);
