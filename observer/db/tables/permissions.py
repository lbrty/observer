from sqlalchemy import Boolean, Column, ForeignKey, Table
from sqlalchemy.dialects.postgresql import UUID

from observer.db.tables import metadata

permissions = Table(
    "permissions",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("can_create", Boolean, default=False, nullable=False),
    Column("can_read", Boolean, default=False, nullable=False),
    Column("can_update", Boolean, default=False, nullable=False),
    Column("can_delete", Boolean, default=False, nullable=False),
    Column("can_read_documents", Boolean, default=False, nullable=False),
    Column("can_read_personal_info", Boolean, default=False, nullable=False),
    Column("user_id", UUID, ForeignKey("users.id"), nullable=False),
    Column("project_id", UUID, ForeignKey("projects.id"), nullable=False),
)
