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
