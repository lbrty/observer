from typing import Literal
from sqlalchemy import (
    Boolean,
    CheckConstraint,
    Column,
    DateTime,
    ForeignKey,
    Index,
    Table,
    Text,
    func,
    text,
)
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, declared_attr, relationship

from observer.db.models import ModelBase, Office
from observer.typing.inspect import unwrap_literal

UserRole = Literal["admin", "staff", "consultant", "guest"]
roles = ", ".join([f"'{role}'" for role in unwrap_literal(UserRole)])


class User(ModelBase):
    __tablename__ = "users"
    __table_args__ = (
        CheckConstraint(f"role IN ({roles})", name="users_role_type_check"),
    )

    email: Mapped[str] = mapped_column("email", Text(), nullable=False, unique=True)
    full_name: Mapped[str] = mapped_column("full_name", Text(), nullable=True)
    password_hash: Mapped[str] = mapped_column("password_hash", Text(), nullable=False)
    role: Mapped[str] = mapped_column("role", Text(), nullable=True)
    is_active: Mapped[str] = mapped_column(
        "is_active",
        Boolean(),
        nullable=True,
        default=True,
    )
    is_confirmed: Mapped[str] = mapped_column(
        "is_confirmed",
        Boolean(),
        nullable=True,
        default=False,
    )
    office_id: Mapped[str] = mapped_column(
        "office_id",
        UUID(as_uuid=True),
        ForeignKey("offices.id", ondelete="SET NULL"),
        nullable=True,
        index=True,
    )
    mfa_enabled: Mapped[str] = mapped_column(
        "mfa_enabled",
        Boolean(),
        nullable=True,
        default=False,
    )
    mfa_encrypted_secret: Mapped[str] = mapped_column(
        "mfa_encrypted_secret",
        Text(),
        nullable=True,
    )
    mfa_encrypted_backup_codes: Mapped[str] = mapped_column(
        "mfa_encrypted_backup_codes",
        Text(),
        nullable=True,
    )

    @declared_attr
    def country(self) -> Mapped[Office]:
        return relationship(
            Office,
            lazy="raise",
            foreign_keys="[offices.id]",
        )
