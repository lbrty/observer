"""world schema

Revision ID: 1565c3d62807
Revises: 7109afee0217
Create Date: 2022-08-10 19:42:15.777905
"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "1565c3d62807"
down_revision = "7109afee0217"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "countries",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.String(length=100), nullable=False),
        sa.Column("code", sa.String(length=4), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ux_countries_name"), "countries", [sa.text("lower(name)")], unique=True)

    op.create_table(
        "states",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.String(length=100), nullable=False),
        sa.Column("code", sa.String(length=10), nullable=False),
        sa.Column("country_id", postgresql.UUID(), nullable=False),
        sa.ForeignKeyConstraint(
            ("country_id",),
            ["users.id"],
            ondelete="cascade",
        ),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ux_states_name"), "states", [sa.text("lower(name)")], unique=True)

    op.create_table(
        "cities",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.String(length=100), nullable=False),
        sa.Column("code", sa.String(length=10), nullable=False),
        sa.Column("state_id", postgresql.UUID(), nullable=False),
        sa.ForeignKeyConstraint(
            ("state_id",),
            ["users.id"],
            ondelete="cascade",
        ),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ux_cities_name"), "cities", [sa.text("lower(name)")], unique=True)


def downgrade():
    op.drop_table("cities")
    op.drop_table("states")
    op.drop_table("countries")
