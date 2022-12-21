"""world

Revision ID: 193c427c31b7
Revises: 7e2a6f5627ca
Create Date: 2022-10-08 22:15:48.379912
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "193c427c31b7"
down_revision = "7e2a6f5627ca"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "countries",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("code", sa.Text(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ux_countries_name"), "countries", [sa.text("lower(name)")], unique=True)
    op.create_index(op.f("ux_countries_code"), "countries", [sa.text("lower(code)")], unique=True)

    op.create_table(
        "states",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("code", sa.Text(), nullable=False),
        sa.Column("country_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.ForeignKeyConstraint(
            ("country_id",),
            ["countries.id"],
            ondelete="CASCADE",
        ),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ux_states_name"), "states", [sa.text("lower(name)")], unique=True)
    op.create_index(op.f("ux_states_code"), "states", [sa.text("lower(code)")], unique=True)
    op.create_index(op.f("ix_states_country_id"), "states", ["country_id"])

    op.create_table(
        "places",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("code", sa.Text(), nullable=False),
        sa.Column("place_type", sa.Text(), nullable=False),
        sa.Column("state_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("country_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.PrimaryKeyConstraint("id"),
        sa.CheckConstraint("place_type IN ('city', 'town', 'village')", name="places_place_types"),
        sa.ForeignKeyConstraint(
            ("state_id",),
            ["states.id"],
            ondelete="CASCADE",
        ),
        sa.ForeignKeyConstraint(
            ("country_id",),
            ["countries.id"],
            ondelete="CASCADE",
        ),
    )

    op.create_index(op.f("ux_places_name"), "places", [sa.text("lower(name)")], unique=True)
    op.create_index(op.f("ux_places_code"), "places", [sa.text("lower(code)")], unique=True)
    op.create_index(op.f("ux_places_place_type"), "places", [sa.text("lower(place_type)")], unique=True)
    op.create_index(op.f("ix_places_state_id"), "places", ["state_id"])
    op.create_index(op.f("ix_places_country_id"), "places", ["country_id"])


def downgrade():
    op.drop_table("places")
    op.drop_table("states")
    op.drop_table("countries")
