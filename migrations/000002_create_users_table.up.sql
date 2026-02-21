CREATE TABLE users (
    id         TEXT        PRIMARY KEY,
    first_name TEXT,
    last_name  TEXT,
    email      VARCHAR(255) NOT NULL,
    phone      VARCHAR(20)  NOT NULL,
    role       VARCHAR(50)  NOT NULL CHECK (role IN ('admin', 'staff', 'consultant', 'guest')),
    is_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT uq_users_phone UNIQUE (phone)
);

CREATE INDEX ix_users_email ON users (email);
CREATE INDEX ix_users_phone ON users (phone);
CREATE INDEX ix_users_role  ON users (role);
