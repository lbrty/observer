from sqlalchemy import UUID, text, Text, TIMESTAMP, func, Index
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.orm import Mapped, mapped_column

from observer.db.models.base import ModelBase


class AuditLog(ModelBase):
    __tablename__ = "audit_logs"

    id: Mapped[UUID] = mapped_column(
        "id",
        UUID(as_uuid=True),
        primary_key=True,
        server_default=text("gen_random_uuid()"),
    )
    ref: Mapped[str] = mapped_column(
        "ref",
        Text(),
        nullable=False,
        index=True,
    )
    data: Mapped[dict | None] = mapped_column(
        "data",
        JSONB(),
        nullable=True,
        default={},
    )
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
