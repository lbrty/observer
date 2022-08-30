"""documents

Revision ID: d97898e0c40b
Revises: 725802c4ebae
Create Date: 2022-08-10 20:29:55.767876
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "d97898e0c40b"
down_revision = "725802c4ebae"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "documents",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("encryption_key", sa.Text(), nullable=True),
        sa.Column("name", sa.String(length=100), nullable=False),
        sa.Column("path", sa.String(length=4), nullable=False),
        sa.Column("person_id", postgresql.UUID(), nullable=False),
        sa.Column("created_at", postgresql.TIMESTAMP(timezone=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
        sa.ForeignKeyConstraint(
            ("person_id",),
            ["displaced_persons.id"],
        ),
    )

    op.create_index(op.f("ix_documents_name"), "documents", [sa.text("lower(name)")])


def downgrade():
    op.drop_table("documents")
