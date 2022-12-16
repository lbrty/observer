from sqlalchemy import Column, ForeignKey, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

countries = Table(
    "countries",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("name", Text(), nullable=False),
    Column("code", Text(), nullable=False),
    Index("ux_countries_name", text("lower(name)"), unique=True),
    Index("ux_countries_code", text("lower(code)"), unique=True),
)

states = Table(
    "states",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("name", Text(), nullable=False),
    Column("code", Text(), nullable=False),
    Column("country_id", UUID(as_uuid=True), ForeignKey("countries.id"), nullable=False),
    Index("ux_states_name", text("lower(name)"), unique=True),
    Index("ux_states_code", text("lower(code)"), unique=True),
    Index("ix_states_country_id", "country_id"),
)

cities = Table(
    "cities",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("name", Text(), nullable=False),
    Column("code", Text(), nullable=False),
    Column("state_id", UUID(as_uuid=True), ForeignKey("states.id"), nullable=False),
    Column("country_id", UUID(as_uuid=True), ForeignKey("countries.id"), nullable=False),
    Index("ux_cities_name", text("lower(name)"), unique=True),
    Index("ux_cities_code", text("lower(code)"), unique=True),
    Index("ix_states_country_id", "country_id"),
    Index("ix_states_state_id", "state_id"),
)
