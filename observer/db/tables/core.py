from sqlalchemy import (
    Boolean,
    CheckConstraint,
    Column,
    ForeignKey,
    Index,
    Table,
    Text,
    text,
)
from sqlalchemy.dialects.postgresql import JSONB, TIMESTAMP, UUID

from observer.db.tables import metadata
from observer.db.util import utcnow

users = Table(
    "users",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("email", Text(), nullable=False),
    Column("full_name", Text(), nullable=True),
    Column("password_hash", Text(), nullable=False),
    Column("role", Text(), nullable=True),
    Column("is_active", Boolean, nullable=True, default=True),
    Column("is_confirmed", Boolean, nullable=True, default=False),
    Column("mfa_enabled", Boolean, nullable=True, default=False),
    Index("ix_users_full_name", text("lower(full_name)")),
    Index("ux_users_email", text("lower(email)"), unique=True),
    Index("ix_users_is_active", "is_active"),
    CheckConstraint("role IN ('admin', 'consultant', 'guest', 'staff')", name="users_role_type_check"),
)

projects = Table(
    "projects",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", Text(), nullable=False, unique=True),
    Column("description", Text(), nullable=True),
    Index("ux_projects_name", text("lower(name)"), unique=True),
)

permissions = Table(
    "permissions",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("can_create", Boolean, default=False, nullable=False),
    Column("can_read", Boolean, default=False, nullable=False),
    Column("can_update", Boolean, default=False, nullable=False),
    Column("can_delete", Boolean, default=False, nullable=False),
    Column("can_read_documents", Boolean, default=False, nullable=False),
    Column("can_read_personal_info", Boolean, default=False, nullable=False),
    Column("user_id", UUID, ForeignKey("users.id"), nullable=False),
    Column("project_id", UUID, ForeignKey("projects.id"), nullable=False),
)

# Audit logs
audit_logs = Table(
    "audit_logs",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("ref", Text(), nullable=False),  # format - origin=<user_id...>;source=services:users;action=create:user;
    Column("data", JSONB(), nullable=True, default={}),
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow),
    Column("expires_at", TIMESTAMP(timezone=True), default=utcnow, nullable=True),
)
