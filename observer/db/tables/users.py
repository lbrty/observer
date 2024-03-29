from sqlalchemy import (
    Boolean,
    CheckConstraint,
    Column,
    DateTime,
    ForeignKey,
    Index,
    Table,
    Text,
    func,
    text,
)
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

users = Table(
    "users",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("email", Text(), nullable=False),
    Column("full_name", Text(), nullable=True),
    Column("password_hash", Text(), nullable=False),
    Column("role", Text(), nullable=True),
    Column("is_active", Boolean(), nullable=True, default=True),
    Column("is_confirmed", Boolean(), nullable=True, default=False),
    Column("office_id", UUID(as_uuid=True), ForeignKey("offices.id", ondelete="SET NULL"), nullable=True),
    Column("mfa_enabled", Boolean(), nullable=True, default=False),
    Column("mfa_encrypted_secret", Text(), nullable=True),
    Column("mfa_encrypted_backup_codes", Text(), nullable=True),
    Index("ix_users_full_name", text("lower(full_name)")),
    Index("ix_users_office_id", text("office_id")),
    Index("ux_users_email", text("lower(email)"), unique=True),
    Index("ix_users_is_active", "is_active"),
    CheckConstraint("role IN ('admin', 'consultant', 'guest', 'staff')", name="users_role_type_check"),
)

password_resets = Table(
    "password_resets",
    metadata,
    Column("code", Text()),
    Column("user_id", UUID(as_uuid=True), ForeignKey("users.id", ondelete="CASCADE"), nullable=False),
    Column("created_at", DateTime(timezone=True), server_default=func.now(), nullable=True),
    Index("ux_password_resets_code", "code", unique=True),
    Index("ix_password_resets_user_id", "user_id"),
)

confirmations = Table(
    "confirmations",
    metadata,
    Column("code", Text()),
    Column("user_id", UUID(as_uuid=True), ForeignKey("users.id", ondelete="CASCADE"), nullable=False),
    Column("expires_at", DateTime(timezone=True), server_default=func.now(), nullable=True),
    Index("ux_confirmations_code", "code", unique=True),
    Index("ix_confirmations_user_id", "user_id"),
)

invites = Table(
    "invites",
    metadata,
    Column("code", Text()),
    Column("user_id", UUID(as_uuid=True), ForeignKey("users.id", ondelete="CASCADE"), nullable=False),
    Column("expires_at", DateTime(timezone=True), server_default=text("now()"), nullable=False),
    Index("ux_invites_code", "code", unique=True),
    Index("ix_invites_user_id", "user_id"),
)
