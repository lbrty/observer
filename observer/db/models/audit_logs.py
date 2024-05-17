from sqlalchemy import UUID, text, Text, TIMESTAMP, func
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.orm import Mapped, mapped_column

from observer.db.models.base import ModelBase


class AuditLog(ModelBase):
    __tablename__ = "audit_logs"

    ref: Mapped[str] = mapped_column(
        "ref",
        Text(),
        nullable=False,
        index=True,
    )
    """Identifier with key-value pairs"""

    data: Mapped[dict | None] = mapped_column(
        "data",
        JSONB(),
        nullable=True,
        default={},
    )
    """Additional metadata for audit logs"""

    created_at: Mapped[TIMESTAMP] = mapped_column(
        "created_at",
        TIMESTAMP(timezone=True),
        default=func.timezone("UTC", func.current_timestamp()),
    )

    expires_at: Mapped[TIMESTAMP] = mapped_column(
        "expires_at",
        TIMESTAMP(timezone=True),
        default=func.timezone("UTC", func.current_timestamp()),
        nullable=True,
        index=True,
    )
    """Expiration datetime"""
