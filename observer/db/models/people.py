from sqlalchemy import (
    CheckConstraint,
    Date,
    ForeignKey,
    Index,
    Text,
    text,
)
from sqlalchemy.dialects.postgresql import ARRAY, JSONB, UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship

from observer.db.choices import age_groups, person_statuses, sex
from observer.db.models import Project, User, TimestampedModel
from observer.db.models.categories import Category
from observer.db.models.offices import Office


class People(TimestampedModel):
    __tablename__ = "people"
    __table_args__ = (
        Index("ix_people_full_name", text("lower(full_name)")),
        Index("ix_people_email", text("lower(email)")),
        CheckConstraint(f"status IN ({', '.join(person_statuses)})", name="status"),
        CheckConstraint(f"sex IN ({', '.join(sex)})", name="sex"),
        CheckConstraint(f"age_group IN ({', '.join(age_groups)})", name="age_group"),
    )

    status: Mapped[str] = mapped_column("status", Text())
    external_id: Mapped[str] = mapped_column("external_id", Text(), nullable=True)
    email: Mapped[str] = mapped_column("email", Text(), nullable=True)
    full_name: Mapped[str] = mapped_column("full_name", Text(), nullable=True)
    birth_date: Mapped[Date] = mapped_column("birth_date", Date, nullable=True)
    sex: Mapped[str] = mapped_column("sex", Text(), nullable=True)
    notes: Mapped[str] = mapped_column("notes", Text(), nullable=True)

    # class PhoneNumbers(encrypted_key: str, items: list[str])
    phone_numbers: Mapped[dict | None] = mapped_column(
        "phone_numbers",
        JSONB(),
        nullable=True,
        default={},
    )

    age_group: Mapped[str] = mapped_column("age_group", Text(), nullable=False)

    project: Mapped[Project] = relationship(Project, back_populates="children")
    project_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("projects.id", ondelete="SET NULL"),
        nullable=True,
        index=True,
    )

    category: Mapped[Project] = relationship(Category, back_populates="children")
    category_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("categories.id", ondelete="SET NULL"),
        nullable=True,
        index=True,
    )

    parent_id = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("people.id", ondelete="CASCADE"),
        nullable=True,
        index=True,
    )

    consultant: Mapped[User] = relationship(User, back_populates="children")
    consultant_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="SET NULL"),
        nullable=True,
        index=True,
    )

    office: Mapped[Project] = relationship(Office, back_populates="children")
    office_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("offices.id", ondelete="SET NULL"),
        nullable=True,
        index=True,
    )

    tags: Mapped[list[str]] = mapped_column(
        "tags",
        ARRAY(Text()),
        nullable=True,
        index=True,
    )
