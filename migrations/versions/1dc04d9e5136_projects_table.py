"""projects_table

Revision ID: 1dc04d9e5136
Revises: 178781a799d5
Create Date: 2022-08-10 19:06:42.550128
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "1dc04d9e5136"
down_revision = "178781a799d5"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "projects",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.String(length=128), nullable=False),
        sa.Column("description", sa.Text(), nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ux_projects_name"), "projects", [sa.text("lower(name)")], unique=True)


def downgrade():
    op.drop_table("projects")
