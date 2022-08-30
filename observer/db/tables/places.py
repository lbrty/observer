from importlib.metadata import metadata

from sqlalchemy import Column, ForeignKey, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import UUID

countries = Table(
    "countries",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", Text(), nullable=False),
    Column("code", Text(), nullable=False),
    Index("ux_countries_name", text("lower(name)"), unique=True),
    Index("ix_countries_code", text("lower(code)")),
)

states = Table(
    "states",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", Text(), nullable=False),
    Column("code", Text(), nullable=False),
    Column("country_id", UUID, ForeignKey("users.id"), nullable=False),
    Index("ux_states_name", text("lower(name)"), unique=True),
    Index("ix_states_code", text("lower(code)")),
)

cities = Table(
    "cities",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", Text(), nullable=False),
    Column("code", Text(), nullable=False),
    Column("state_id", UUID, ForeignKey("users.id"), nullable=False),
    Index("ux_cities_name", text("lower(name)"), unique=True),
    Index("ix_cities_code", text("lower(code)")),
)
