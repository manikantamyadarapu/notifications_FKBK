-- Minimal schema aligned with the coding assessment requirements.
-- Tables used by this microservice:
-- - tamper_events (event ingestion / storage)
-- - tamper_code_desc (mandatory tamper_code -> description mapping)
-- - escalation_notifications (tamper alert notifications)
-- - meters (meter_id validation / enrichment)

CREATE TABLE IF NOT EXISTS meters (
  meter_id TEXT PRIMARY KEY,
  -- optional enrichment fields (extend as needed)
  meter_name TEXT,
  location TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tamper_code_desc (
  tamper_code INTEGER PRIMARY KEY,
  tamper_desc TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tamper_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  meter_id TEXT NOT NULL,
  tamper_code INTEGER NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  processed BOOLEAN NOT NULL DEFAULT FALSE,
  FOREIGN KEY (meter_id) REFERENCES meters(meter_id)
);

CREATE INDEX IF NOT EXISTS idx_tamper_events_processed ON tamper_events(processed);
CREATE INDEX IF NOT EXISTS idx_tamper_events_meter_time ON tamper_events(meter_id, timestamp);

CREATE TABLE IF NOT EXISTS escalation_notifications (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  meter_id TEXT NOT NULL,
  tamper_code INTEGER NOT NULL,
  tamper_description TEXT NOT NULL,
  message TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (meter_id) REFERENCES meters(meter_id)
);

CREATE INDEX IF NOT EXISTS idx_escalation_notifications_meter_time
  ON escalation_notifications(meter_id, timestamp);

CREATE INDEX IF NOT EXISTS idx_escalation_notifications_code_time
  ON escalation_notifications(tamper_code, timestamp);

