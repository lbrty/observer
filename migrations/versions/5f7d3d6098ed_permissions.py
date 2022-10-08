"""permissions

Revision ID: 5f7d3d6098ed
Revises: d75cff45a8cd
Create Date: 2022-10-08 22:14:29.378657
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "5f7d3d6098ed"
down_revision = "d75cff45a8cd"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "permissions",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("can_create", sa.Boolean(), nullable=False),
        sa.Column("can_read", sa.Boolean(), nullable=False),
        sa.Column("can_update", sa.Boolean(), nullable=False),
        sa.Column("can_delete", sa.Boolean(), nullable=False),
        sa.Column("can_read_documents", sa.Boolean(), nullable=False),
        sa.Column("can_read_personal_info", sa.Boolean(), nullable=False),
        sa.Column("user_id", postgresql.UUID(), nullable=False),
        sa.Column("project_id", postgresql.UUID(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
        sa.ForeignKeyConstraint(
            ("project_id",),
            ["projects.id"],
        ),
        sa.ForeignKeyConstraint(
            ("user_id",),
            ["users.id"],
        ),
    )


def downgrade():
    op.drop_table("permissions")