"""offices

Revision ID: de0f408dabcf
Revises: 2e589e3665c8
Create Date: 2023-01-31 18:35:45.256677
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "de0f408dabcf"
down_revision = "2e589e3665c8"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "offices",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("place_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
        sa.ForeignKeyConstraint(
            ("place_id",),
            ["places.id"],
            ondelete="SET NULL",
        ),
    )

    op.create_index(op.f("ix_offices_place_id"), "offices", ["place_id"])


def downgrade():
    op.drop_table("offices")
