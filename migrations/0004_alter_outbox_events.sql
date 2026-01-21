ALTER TABLE outbox_events
  ALTER COLUMN aggregate_id TYPE BIGINT
  USING aggregate_id::bigint;
