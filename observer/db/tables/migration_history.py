from sqlalchemy import Column, Date, DateTime, ForeignKey, Index, Table, func, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

migration_history = Table(
    "migration_history",
    metadata,
    Column(
        "id",
        UUID(as_uuid=True),
        nullable=False,
        server_default=text("gen_random_uuid()"),
    ),
    Column(
        "idp_id",
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="CASCADE"),
        nullable=False,
    ),
    Column("migration_date", Date(), nullable=True),
    Column(
        "project_id",
        UUID(as_uuid=True),
        ForeignKey("projects.id", ondelete="SET NULL"),
        nullable=True,
    ),
    Column(
        "from_place_id",
        UUID(as_uuid=True),
        ForeignKey("places.id", ondelete="SET NULL"),
        nullable=True,
    ),
    Column(
        "current_place_id",
        UUID(as_uuid=True),
        ForeignKey("places.id", ondelete="SET NULL"),
        nullable=True,
    ),
    Column(
        "created_at",
        DateTime(timezone=True),
        server_default=func.now(),
        nullable=True,
    ),
    Index("ix_migration_history_idp_id", "idp_id"),
    Index("ix_migration_history_project_id", "project_id"),
    Index("ix_migration_history_migration_date", "migration_date"),
    Index("ix_migration_history_from_place_id", "from_place_id"),
    Index("ix_migration_history_current_place_id", "current_place_id"),
)
