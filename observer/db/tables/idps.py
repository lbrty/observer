from sqlalchemy import (
    Column,
    Index,
    ForeignKey,
    Table,
    Text,
    text,
)
from sqlalchemy.dialects.postgresql import ARRAY, DATE, UUID, TIMESTAMP

from observer.db.tables import metadata
from observer.db.util import utcnow

vulnerability_categories = Table(
    "vulnerability_categories",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", Text(), nullable=False),
    Index("ux_categories_name", text("lower(name)"), unique=True),
)

displaced_persons = Table(
    "displaced_persons",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("status", Text()),
    Column("external_id", Text()),
    Column("reference_id", Text()),
    Column("email", Text(), nullable=True),
    Column("full_name", Text(), nullable=False),
    Column("birth_date", DATE, nullable=True),
    Column("notes", Text(), nullable=True),
    Column("phone_number", Text(), nullable=True),
    Column("phone_number_additional", Text(), nullable=True),
    Column("migration_date", DATE, nullable=True),
    # Location info
    Column("from_city_id", UUID, ForeignKey("cities.id"), nullable=True),
    Column("from_state_id", UUID, ForeignKey("states.id"), nullable=True),
    Column("current_city_id", UUID, ForeignKey("cities.id"), nullable=True),
    Column("current_state_id", UUID, ForeignKey("states.id"), nullable=True),
    Column("project_id", UUID, ForeignKey("projects.id"), nullable=True),
    Column("category_id", UUID, ForeignKey("vulnerability_categories.id"), nullable=True),
    # User's id who registered
    Column("creator_id", UUID, ForeignKey("users.id"), nullable=True),
    Column("tags", ARRAY(Text()), nullable=True),
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow),
    Column("updated_at", TIMESTAMP(timezone=True), default=utcnow, onupdate=utcnow),
)

documents = Table(
    "documents",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("encryption_key", Text()),
    Column("name", Text(), nullable=False),
    Column("path", Text(), nullable=False),
    Column("person_id", UUID, ForeignKey("displaced_persons.id"), nullable=False),
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow),
    Index("ux_documents_name", text("lower(name)")),
)
