"""categories

Revision ID: 7e2a6f5627ca
Revises: b9a0f0c6205d
Create Date: 2022-10-08 22:15:21.249910
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "7e2a6f5627ca"
down_revision = "b9a0f0c6205d"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "categories",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.Text(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ux_categories_name"), "categories", [sa.text("lower(name)")], unique=True)


def downgrade():
    op.drop_table("categories")
