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
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("description", sa.Text(), nullable=False),
        sa.Column("owner_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
        sa.ForeignKeyConstraint(
            ("owner_id",),
            ["users.id"],
            ondelete="SET NULL",
        ),
    )

    op.create_index(op.f("ix_projects_name"), "projects", [sa.text("lower(name)")])
    op.create_index(op.f("ix_projects_owner_id"), "projects", [sa.text("owner_id")])


def downgrade():
    op.drop_table("projects")
