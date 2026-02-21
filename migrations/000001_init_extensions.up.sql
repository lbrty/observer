-- Enable useful PostgreSQL extensions

-- Case-insensitive text type (useful for emails, usernames)
CREATE EXTENSION IF NOT EXISTS citext;

-- Trigram matching for fuzzy text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Remove accents from text (useful for Cyrillic/Latin search)
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Cryptographic functions
CREATE EXTENSION IF NOT EXISTS pgcrypto;
