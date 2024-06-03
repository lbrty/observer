# from sqlalchemy import (
#     CheckConstraint,
#     Column,
#     DateTime,
#     ForeignKey,
#     Index,
#     Table,
#     Text,
#     func,
#     text,
# )
# from sqlalchemy.dialects.postgresql import UUID
#
# from observer.db import metadata
#
# support_records = Table(
#     "support_records",
#     metadata,
#     Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
#     Column("description", Text(), nullable=True),
#     Column("type", Text(), nullable=False),
#     Column("consultant_id", UUID(as_uuid=True), nullable=False),
#     Column("record_for", Text(), nullable=False),
#     Column("owner_id", UUID(as_uuid=True), nullable=False),
#     Column(
#         "project_id",
#         UUID(as_uuid=True),
#         ForeignKey("projects.id", ondelete="SET NULL"),
#         nullable=True,
#     ),
#     Column("created_at", DateTime(timezone=True), server_default=func.now(), nullable=True),
#     Index("ix_support_records_type", "type"),
#     Index("ix_support_records_description", "description"),
#     Index("ix_support_records_consultant_id", "consultant_id"),
#     Index("ix_support_records_age_group", "age_group"),
#     Index("ix_support_records_owner_id", "owner_id"),
#     Index("ix_support_records_project_id", "project_id"),
#     CheckConstraint("type IN ('humanitarian', 'legal', 'medical', 'general')", name="support_records_types"),
#     CheckConstraint("record_for IN ('person', 'pet')", name="support_records_record_for"),
#     CheckConstraint(
#         (
#             "age_group IN ("
#             "'infant', 'toddler', 'pre_school', "
#             "'middle_childhood', 'young_teen', "
#             "'teenager', 'young_adult', 'early_adult', "
#             "'middle_aged_adult', 'old_adult'"
#             ")"
#         ),
#         name="support_records_age_group",
#     ),
# )
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
