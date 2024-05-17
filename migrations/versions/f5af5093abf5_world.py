"""world

Revision ID: f5af5093abf5
Revises: c685c35f23b2
Create Date: 2024-05-11 16:02:15.002987
"""

from typing import Sequence, Union

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "f5af5093abf5"
down_revision: Union[str, None] = "c685c35f23b2"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "countries",
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("code", sa.Text(), nullable=False),
        sa.Column(
            "id",
            sa.UUID(),
            server_default=sa.text("gen_random_uuid()"),
            nullable=False,
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_countries")),
        sa.UniqueConstraint("code", name=op.f("uq_countries_code_key")),
    )
    op.create_index(
        op.f("ix_countries_name"),
        "countries",
        ["name"],
    )
    op.create_table(
        "states",
        sa.Column("country_id", sa.UUID(), nullable=False),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("code", sa.Text(), nullable=False),
        sa.Column(
            "id",
            sa.UUID(),
            server_default=sa.text("gen_random_uuid()"),
            nullable=False,
        ),
        sa.ForeignKeyConstraint(
            ["country_id"],
            ["countries.id"],
            name=op.f("fk_states_country_id_countries"),
            ondelete="CASCADE",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_states")),
        sa.UniqueConstraint("code", name=op.f("uq_states_code_key")),
    )
    op.create_index(
        op.f("ix_states_country_id"),
        "states",
        ["country_id"],
    )
    op.create_index(
        op.f("ix_states_name"),
        "states",
        ["name"],
    )
    op.create_table(
        "places",
        sa.Column("state_id", sa.UUID(), nullable=False),
        sa.Column("country_id", sa.UUID(), nullable=False),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("code", sa.Text(), nullable=False),
        sa.Column(
            "id",
            sa.UUID(),
            server_default=sa.text("gen_random_uuid()"),
            nullable=False,
        ),
        sa.ForeignKeyConstraint(
            ["country_id"],
            ["countries.id"],
            name=op.f("fk_places_country_id_countries"),
            ondelete="CASCADE",
        ),
        sa.ForeignKeyConstraint(
            ["state_id"],
            ["states.id"],
            name=op.f("fk_places_state_id_states"),
            ondelete="CASCADE",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_places")),
        sa.UniqueConstraint("code", name=op.f("uq_places_code_key")),
    )
    op.create_index(
        op.f("ix_places_country_id"),
        "places",
        ["country_id"],
    )
    op.create_index(
        op.f("ix_places_name"),
        "places",
        ["name"],
    )
    op.create_index(
        op.f("ix_places_state_id"),
        "places",
        ["state_id"],
    )


def downgrade() -> None:
    op.drop_index(op.f("ix_places_state_id"), table_name="places")
    op.drop_index(op.f("ix_places_name"), table_name="places")
    op.drop_index(op.f("ix_places_country_id"), table_name="places")
    op.drop_index(op.f("ix_states_name"), table_name="states")
    op.drop_index(op.f("ix_states_country_id"), table_name="states")
    op.drop_index(op.f("ix_countries_name"), table_name="countries")
    op.drop_table("places")
    op.drop_table("states")
    op.drop_table("countries")
