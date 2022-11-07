from sqlalchemy import Boolean, CheckConstraint, Column, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db.tables import metadata

users = Table(
    "users",
    metadata,
    Column("id", UUID(), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("ref_id", Text(), primary_key=True),
    Column("email", Text(), nullable=False),
    Column("full_name", Text(), nullable=True),
    Column("password_hash", Text(), nullable=False),
    Column("role", Text(), nullable=True),
    Column("is_active", Boolean(), nullable=True, default=True),
    Column("is_confirmed", Boolean(), nullable=True, default=False),
    Column("mfa_enabled", Boolean(), nullable=True, default=False),
    Column("mfa_encrypted_secret", Text(), nullable=True),
    Column("mfa_encrypted_backup_codes", Text(), nullable=True),
    Index("ix_users_full_name", text("lower(full_name)")),
    Index("ux_users_email", text("lower(email)"), unique=True),
    Index("ix_users_is_active", "is_active"),
    CheckConstraint("role IN ('admin', 'consultant', 'guest', 'staff')", name="users_role_type_check"),
)
