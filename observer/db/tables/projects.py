from sqlalchemy import Column, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

projects = Table(
    "projects",
    metadata,
    Column("id", UUID(), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("name", Text(), nullable=False),
    Column("description", Text(), nullable=True),
    Index("ix_projects_name", text("lower(name)")),
)
