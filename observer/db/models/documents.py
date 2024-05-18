from sqlalchemy import Text, Integer, UUID, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from observer.db.models import ModelBase, User, Project


class Document(ModelBase):
    __tablename__ = "documents"

    encryption_key: Mapped[str | None] = mapped_column(
        "encryption_key",
        Text(),
        nullable=True,
    )

    name: Mapped[str] = mapped_column(
        "name",
        Text(),
        nullable=False,
        index=True,
    )

    path: Mapped[str] = mapped_column(
        "path",
        Text(),
        nullable=False,
    )

    mime: Mapped[str] = mapped_column(
        "mime",
        Text(),
        nullable=False,
    )

    size: Mapped[float] = mapped_column(
        "size",
        Integer(),
        nullable=False,
    )

    owner_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    owner: Mapped[User] = relationship(
        User,
        back_populates="parent",
    )

    project_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("projects.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    project: Mapped[Project] = relationship(
        Project,
        back_populates="parent",
    )
