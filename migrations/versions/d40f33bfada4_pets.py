"""pets

Revision ID: d40f33bfada4
Revises: d97898e0c40b
Create Date: 2022-10-03 23:41:00.829837
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

from observer.db.util import utcnow

# revision identifiers, used by Alembic.
revision = "d40f33bfada4"
down_revision = "d97898e0c40b"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "pets",
        sa.Column("id", postgresql.UUID, primary_key=True),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("notes", sa.Text(), nullable=True),
        sa.Column("status", sa.Text(), nullable=False),
        sa.Column("registration_id", sa.Text(), nullable=True),  # registration document id
        sa.Column("created_at", postgresql.TIMESTAMP(timezone=True), default=utcnow),
        sa.Column("owner_id", postgresql.UUID, nullable=True),
        sa.ForeignKeyConstraint(
            ("owner_id",),
            ["users.id"],
            ondelete="SET NULL",
        ),
    )

    op.create_index(op.f("ix_pets_name"), "pets", [sa.text("lower(name)")])
    op.create_index(op.f("ix_pets_status"), "pets", ["status"])
    op.create_index(op.f("ix_pets_registration_id"), "pets", ["registration_id"])


def downgrade():
    op.drop_table("pets")
