from sqlalchemy import Column, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db.tables import metadata

projects = Table(
    "projects",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", Text(), nullable=False, unique=True),
    Column("description", Text(), nullable=True),
    Index("ix_projects_name", text("lower(name)")),
)