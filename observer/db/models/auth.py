from sqlalchemy import Text, UUID, ForeignKey, func, TIMESTAMP
from sqlalchemy.orm import Mapped, mapped_column

from observer.db.models import ModelBase


class PasswordReset(ModelBase):
    __tablename__ = "password_resets"
    name: Mapped[str] = mapped_column(
        "code",
        Text(),
        nullable=False,
        index=True,
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


# confirmations = Table(
#     "confirmations",
#     metadata,
#     Column("code", Text()),
#     Column("user_id", UUID(as_uuid=True), ForeignKey("users.id", ondelete="CASCADE"), nullable=False),
#     Column("expires_at", DateTime(timezone=True), server_default=func.now(), nullable=True),
#     Index("ux_confirmations_code", "code", unique=True),
#     Index("ix_confirmations_user_id", "user_id"),
# )
#
# invites = Table(
#     "invites",
#     metadata,
#     Column("code", Text()),
#     Column("user_id", UUID(as_uuid=True), ForeignKey("users.id", ondelete="CASCADE"), nullable=False),
#     Column("expires_at", DateTime(timezone=True), server_default=text("now()"), nullable=False),
#     Index("ux_invites_code", "code", unique=True),
#     Index("ix_invites_user_id", "user_id"),
# )
