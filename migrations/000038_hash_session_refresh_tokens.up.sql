CREATE EXTENSION IF NOT EXISTS pgcrypto;

UPDATE sessions
SET refresh_token = encode(digest(refresh_token, 'sha256'), 'hex');
