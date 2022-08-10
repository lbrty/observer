"""categories

Revision ID: b4849b0f8d04
Revises: 1565c3d62807
Create Date: 2022-08-10 19:53:53.514571
"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "b4849b0f8d04"
down_revision = "1565c3d62807"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "categories",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.String(length=64), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ux_categories_name"), "categories", [sa.text("lower(name)")], unique=True)


def downgrade():
    op.drop_table("categories")
