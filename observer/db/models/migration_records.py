from sqlalchemy import TIMESTAMP, UUID, ForeignKey, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from observer.db.models import TimestampedModel
from observer.db.models.people import People
from observer.db.models.projects import Project
from observer.db.models.world import Place


class MigrationRecord(TimestampedModel):
    __tablename__ = "migration_records"

    notes: Mapped[str] = mapped_column(
        "notes",
        Text(),
        nullable=True,
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
        nullable=False,
        index=True,
    )

    current_place: Mapped[Place] = relationship(Place, back_populates="children")
    current_place_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("places.id", ondelete="SET NULL"),
        nullable=False,
        index=True,
    )

    from_place: Mapped[Place] = relationship(Place, back_populates="children")
    from_place_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("places.id", ondelete="SET NULL"),
        nullable=False,
        index=True,
    )
