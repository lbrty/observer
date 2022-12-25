"""pets

Revision ID: 98cf562f08e2
Revises: e07cd90eaa03
Create Date: 2022-10-08 22:16:41.157524
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "98cf562f08e2"
down_revision = "e07cd90eaa03"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "pets",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("notes", sa.Text(), nullable=True),
        sa.Column("status", sa.Text(), nullable=False),
        sa.Column("registration_id", sa.Text(), nullable=True),  # registration document id
        sa.Column("owner_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("created_at", sa.DateTime(timezone=True), server_default=sa.text("now()"), nullable=True),
        sa.PrimaryKeyConstraint("id"),
        sa.ForeignKeyConstraint(
            ("owner_id",),
            ["people.id"],
            ondelete="SET NULL",
        ),
    )

    op.create_index(op.f("ix_pets_name"), "pets", [sa.text("lower(name)")])
    op.create_index(op.f("ix_pets_status"), "pets", ["status"])
    op.create_index(op.f("ix_pets_owner_id"), "pets", ["owner_id"])
    op.create_index(op.f("ix_pets_registration_id"), "pets", ["registration_id"])


def downgrade():
    op.drop_table("pets")
