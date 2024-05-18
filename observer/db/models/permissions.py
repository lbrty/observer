from sqlalchemy import Boolean, Text, UUID, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column

from observer.db.models import ModelBase, User, Project


class Permission(ModelBase):
    __tablename__ = "permissions"

    note: Mapped[str] = mapped_column(
        "note",
        Text(),
        nullable=False,
        unique=True,
    )

    can_create: Mapped[bool] = mapped_column(
        "can_create",
        Boolean(),
        default=False,
        nullable=False,
    )

    can_read: Mapped[bool] = mapped_column(
        "can_read",
        Boolean(),
        default=False,
        nullable=False,
    )

    can_update: Mapped[bool] = mapped_column(
        "can_update",
        Boolean(),
        default=False,
        nullable=False,
    )

    can_delete: Mapped[bool] = mapped_column(
        "can_delete",
        Boolean(),
        default=False,
        nullable=False,
    )

    can_read_documents: Mapped[bool] = mapped_column(
        "can_read_documents",
        Boolean(),
        default=False,
        nullable=False,
    )

    can_read_personal_info: Mapped[bool] = mapped_column(
        "can_read_personal_info",
        Boolean(),
        default=False,
        nullable=False,
    )

    can_invite_members: Mapped[bool] = mapped_column(
        "can_invite_members",
        Boolean(),
        default=False,
        nullable=False,
    )

    user_id: Mapped[bool] = mapped_column(
        "user_id",
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    project_id: Mapped[bool] = mapped_column(
        "project_id",
        UUID(as_uuid=True),
        ForeignKey("projects.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )
