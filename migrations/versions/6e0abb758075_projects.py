"""projects

Revision ID: 6e0abb758075
Revises: 47cc3e1dfb71
Create Date: 2022-10-08 22:09:43.300945
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "6e0abb758075"
down_revision = "47cc3e1dfb71"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "projects",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("description", sa.Text(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ix_projects_name"), "projects", [sa.text("lower(name)")])


def downgrade():
    op.drop_table("projects")