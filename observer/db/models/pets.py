# from sqlalchemy import Column, DateTime, ForeignKey, Index, Table, Text, func, text
# from sqlalchemy.dialects.postgresql import UUID
#
# from observer.db import metadata
#
# pets = Table(
#     "pets",
#     metadata,
#     Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
#     Column("name", Text(), nullable=False),
#     Column("notes", Text(), nullable=True),
#     Column("status", Text(), nullable=False),
#     Column("registration_id", Text(), nullable=True),  # registration document id
#     Column("created_at", DateTime(timezone=True), server_default=func.now(), nullable=True),
#     Column("owner_id", UUID(as_uuid=True), ForeignKey("people.id", ondelete="SET NULL"), nullable=True),
#     Column("project_id", UUID(as_uuid=True), ForeignKey("projects.id", ondelete="SET NULL"), nullable=True),
#     Index("ix_pets_name", text("lower(name)")),
#     Index("ix_pets_status", text("status")),
#     Index("ix_pets_registration_id", text("status")),
#     Index("ix_pets_owner_id", text("owner_id")),
#     Index("ix_pets_project_id", text("project_id")),
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


class Pet(ModelBase):
    __tablename__ = "pets"
    __table_args__ = (
        Index("ix_name_pets", text("lower(name)")),
        CheckConstraint(f"status IN ({statuses})", name="pets_status_check"),
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

    registration_id: Mapped[str] = mapped_column(
        "registration_id",
        Text(),
        index=True,
    )

    created_at: Mapped[TIMESTAMP] = mapped_column(
        "created_at",
        TIMESTAMP(timezone=True),
        default=func.timezone("UTC", func.current_timestamp()),
        nullable=False,
    )

    owner_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="SET NULL"),
        nullable=False,
        index=True,
    )

    project_id: Mapped[bool] = mapped_column(
        "project_id",
        UUID(as_uuid=True),
        ForeignKey("projects.id", ondelete="SET NULL"),
        nullable=False,
        index=True,
    )
