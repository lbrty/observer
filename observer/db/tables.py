from sqlalchemy import (
    Boolean,
    Column,
    Index,
    MetaData,
    String,
    Table,
    text,
)
from sqlalchemy.dialects.postgresql import UUID, TIMESTAMP
from observer.db.util import utcnow


metadata = MetaData(
    {
        "ix": "ix_%(column_0_label)s",
        "ux": "ux_%(table_name)s_%(column_0_name)s",
        "ck": "ck_%(table_name)s_%(constraint_name)s",
        "fk": "fk_%(table_name)s_%(column_0_name)s_%(referred_table_name)s",
        "pk": "pk_%(table_name)s",
    }
)

users = Table(
    "users",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("email", String(255), nullable=False),
    Column("full_name", String(128), nullable=True),
    Column("password_hash", String(512), nullable=False),
    Column("role", String(length=20), nullable=True),
    Column("is_active", Boolean, nullable=True, default=True),
    Column("is_confirmed", Boolean, nullable=True, default=False),
    Column("mfa_enabled", Boolean, nullable=True, default=False),
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow),
    Column("updated_at", TIMESTAMP(timezone=True), default=utcnow, onupdate=utcnow),
    Index("ix_users_full_name", text("lower(full_name)")),
    Index("ix_users_email", text("lower(email)"), unique=True),
    Index("ix_users_is_active", "is_active"),
)

displaced_persons = Table("displaced_persons", metadata)
