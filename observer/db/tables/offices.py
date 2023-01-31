from sqlalchemy import Column, Table, Text, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

offices = Table(
    "offices",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("name", Text(), nullable=False),
)
