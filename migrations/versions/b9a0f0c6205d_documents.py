"""documents

Revision ID: b9a0f0c6205d
Revises: 5f7d3d6098ed
Create Date: 2022-10-08 22:14:59.997090
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "b9a0f0c6205d"
down_revision = "5f7d3d6098ed"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "documents",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("encryption_key", sa.Text(), nullable=True),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("path", sa.Text(), nullable=False),
        sa.Column("owner_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("created_at", sa.DateTime(timezone=True), server_default=sa.text("now()"), nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ix_documents_name"), "documents", [sa.text("lower(name)")])
    op.create_index(op.f("ix_documents_owner_id"), "documents", ["owner_id"])


def downgrade():
    op.drop_table("documents")
