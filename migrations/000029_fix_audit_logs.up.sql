-- Make user_id nullable and change cascade to SET NULL so deleted users
-- don't wipe their audit trail (GDPR Article 30).
ALTER TABLE audit_logs
    ALTER COLUMN user_id DROP NOT NULL,
    DROP CONSTRAINT audit_logs_user_id_fkey,
    ADD CONSTRAINT audit_logs_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;

-- Index for per-entity history lookups.
CREATE INDEX ix_audit_logs_entity ON audit_logs (entity_type, entity_id, created_at DESC);
