from sqlalchemy import TIMESTAMP, UUID, CheckConstraint, ForeignKey, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from observer.db.choices import support_types
from observer.db.models import TimestampedModel
from observer.db.models.people import People
from observer.db.models.projects import Project
from observer.db.models.users import User


class SupportRecord(TimestampedModel):
    __tablename__ = "support_records"
    __table_args__ = (CheckConstraint(f"type IN ({', '.join(support_types)})", name="type"),)

    notes: Mapped[str] = mapped_column(
        "notes",
        Text(),
        nullable=True,
    )

    type: Mapped[str] = mapped_column(
        "type",
        Text(),
        nullable=False,
    )

    owner_id: Mapped[str] = mapped_column(
        "owner_id",
        Text(),
        nullable=False,
        index=True,
    )

    consultant: Mapped[User] = relationship(User, back_populates="children")
    consultant_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="SET NULL"),
        nullable=True,
        index=True,
    )

    migration_date: Mapped[TIMESTAMP] = mapped_column(
        "migration_date",
        TIMESTAMP(timezone=True),
        nullable=True,
    )

    person: Mapped[People] = relationship(People, back_populates="children")
    person_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("people.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    project: Mapped[Project] = relationship(Project, back_populates="children")
    project_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("projects.id", ondelete="SET NULL"),
        nullable=True,
        index=True,
    )
