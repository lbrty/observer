"""permissions_table

Revision ID: 7109afee0217
Revises: 1dc04d9e5136
Create Date: 2022-08-10 19:21:58.225226
"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "7109afee0217"
down_revision = "1dc04d9e5136"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "permissions",
        sa.Column("id", postgresql.UUID(), nullable=False),
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
