ALTER TABLE outbox_events
  ALTER COLUMN aggregate_id TYPE UUID
  USING gen_random_uuid();
