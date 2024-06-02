# from sqlalchemy import Column, Date, DateTime, ForeignKey, Index, Table, func, text
# from sqlalchemy.dialects.postgresql import UUID
#
# from observer.db import metadata
#
# migration_history = Table(
#     "migration_history",
#     metadata,
#     Column(
#         "id",
#         UUID(as_uuid=True),
#         nullable=False,
#         server_default=text("gen_random_uuid()"),
#     ),
#     Column(
#         "person_id",
#         UUID(as_uuid=True),
#         ForeignKey("people.id", ondelete="CASCADE"),
#         nullable=False,
#     ),
#     Column("migration_date", Date(), nullable=True),
#     Column(
#         "project_id",
#         UUID(as_uuid=True),
#         ForeignKey("projects.id", ondelete="SET NULL"),
#         nullable=True,
#     ),
#     Column(
#         "from_place_id",
#         UUID(as_uuid=True),
#         ForeignKey("places.id", ondelete="SET NULL"),
#         nullable=True,
#     ),
#     Column(
#         "current_place_id",
#         UUID(as_uuid=True),
#         ForeignKey("places.id", ondelete="SET NULL"),
#         nullable=True,
#     ),
#     Column(
#         "created_at",
#         DateTime(timezone=True),
#         server_default=func.now(),
#         nullable=True,
#     ),
#     Index("ix_migration_history_person_id", "person_id"),
#     Index("ix_migration_history_project_id", "project_id"),
#     Index("ix_migration_history_migration_date", "migration_date"),
#     Index("ix_migration_history_from_place_id", "from_place_id"),
#     Index("ix_migration_history_current_place_id", "current_place_id"),
# )
from sqlalchemy import (
    Text,
    Index,
    text,
    CheckConstraint,
    TIMESTAMP,
    func,
    UUID,
    ForeignKey,
)
from sqlalchemy.orm import Mapped, mapped_column

from observer.common.reflect.inspect import unwrap_enum
from observer.common.types import PetStatus
from observer.db.models import ModelBase

statuses = list(unwrap_enum(PetStatus).keys())


class MigrationRecord(ModelBase):
    __tablename__ = "migration_records"
    __table_args__ = (
        Index("ix_name_pets", text("lower(name)")),
        CheckConstraint(f"status IN ({statuses})", name="ck_pets_status"),
    )

    name: Mapped[str] = mapped_column(
        "name",
        Text(),
    )

    notes: Mapped[str] = mapped_column(
        "notes",
        Text(),
        nullable=True,
    )

    status: Mapped[str] = mapped_column(
        "status",
        Text(),
        index=True,
    )
