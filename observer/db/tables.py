from sqlalchemy import (
    Boolean,
    CheckConstraint,
    Column,
    Index,
    UniqueConstraint,
    ForeignKey,
    MetaData,
    String,
    Table,
    Text,
    text,
)
from sqlalchemy.dialects.postgresql import ARRAY, DATE, UUID, TIMESTAMP
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
    Index("ix_users_full_name", text("lower(full_name)")),
    Index("ux_users_email", text("lower(email)"), unique=True),
    Index("ix_users_is_active", "is_active"),
    CheckConstraint("role IN ('admin', 'consultant', 'guest', 'staff')", name="users_role_type_check"),
)

projects = Table(
    "projects",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", String(128), nullable=False, unique=True),
    Column("description", Text(), nullable=True),
    Index("ix_projects_name", text("lower(name)"), unique=True),
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

# displaced_persons = Table(
#     "displaced_persons",
#     metadata,
#     Column("id", UUID, primary_key=True),
#     Column("status", String(20)),
#     Column("external_id", String(128)),
#     Column("reference_id", String(128)),
#     Column("email", String(255), nullable=False),
#     Column("full_name", String(128), nullable=True),
#     Column("birth_date", DATE, nullable=True),
#     Column("notes", Text(), nullable=True),
#     Column("phone_number", String(20), nullable=True),
#     Column("phone_number_additional", String(20), nullable=True),
#     Column("migration_date", DATE, nullable=True),
#
#     # Location info
#     Column("from_city_id", UUID, nullable=True),
#     Column("from_state_id", UUID, nullable=True),
#     Column("current_city_id", UUID, nullable=True),
#     Column("current_state_id", UUID, nullable=True),
#
#     Column("project_id", UUID, nullable=True),
#     Column("category_id", UUID, nullable=True),
#
#     # User's id who registered
#     Column("creator_id", UUID, nullable=True),
#
#     Column("tags", ARRAY(String(20)), nullable=True),
#     Column("created_at", TIMESTAMP(timezone=True), default=utcnow),
#     Column("updated_at", TIMESTAMP(timezone=True), default=utcnow, onupdate=utcnow),
# )

countries = Table(
    "countries",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", String(255), nullable=False),
    Column("code", String(4), nullable=False),
    Index("ix_countries_name", text("lower(name)")),
    Index("ix_countries_code", text("lower(code)")),
    UniqueConstraint("name"),
)

states = Table(
    "states",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", String(255), nullable=False),
    Column("code", String(10), nullable=False),
    Column("country_id", UUID, ForeignKey("users.id"), nullable=False),
    Index("ix_states_name", text("lower(name)")),
    Index("ix_states_code", text("lower(code)")),
    UniqueConstraint("name"),
)

cities = Table(
    "cities",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", String(255), nullable=False),
    Column("code", String(10), nullable=False),
    Column("state_id", UUID, ForeignKey("users.id"), nullable=False),
    Index("ix_cities_name", text("lower(name)")),
    Index("ix_cities_code", text("lower(code)")),
    UniqueConstraint("name"),
)
