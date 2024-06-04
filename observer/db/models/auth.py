from sqlalchemy import Text, UUID, ForeignKey, func, TIMESTAMP
from sqlalchemy.orm import Mapped, mapped_column

from observer.db.models import ModelBase


class PasswordReset(ModelBase):
    __tablename__ = "password_resets"
    code: Mapped[str] = mapped_column(
        "code",
        Text(),
        nullable=False,
        unique=True,
    )

    user_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    created_at: Mapped[TIMESTAMP] = mapped_column(
        "created_at",
        TIMESTAMP(timezone=True),
        default=func.timezone("UTC", func.current_timestamp()),
        nullable=False,
    )


class Confirmation(ModelBase):
    __tablename__ = "confirmations"
    code: Mapped[str] = mapped_column(
        "code",
        Text(),
        nullable=False,
        unique=True,
    )

    user_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    expires_at: Mapped[TIMESTAMP] = mapped_column(
        "expires_at",
        TIMESTAMP(timezone=True),
        default=func.timezone("UTC", func.current_timestamp()),
        nullable=True,
    )


class Invite(ModelBase):
    __tablename__ = "invites"
    code: Mapped[str] = mapped_column(
        "code",
        Text(),
        nullable=False,
        unique=True,
    )

    user_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    expires_at: Mapped[TIMESTAMP] = mapped_column(
        "expires_at",
        TIMESTAMP(timezone=True),
        default=func.timezone("UTC", func.current_timestamp()),
        nullable=True,
    )
