"""initial

Revision ID: 178781a799d5
Revises:
Create Date: 2022-08-07 22:21:29.485044
"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "178781a799d5"
down_revision = None
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "users",
        sa.Column("id", postgresql.UUID(), nullable=False),
        sa.Column("email", sa.String(length=255), nullable=False),
        sa.Column("full_name", sa.String(length=128), nullable=True),
        sa.Column("password_hash", sa.String(length=512), nullable=False),
        sa.Column("is_active", sa.Boolean(), nullable=True, server_default="1"),
        sa.Column("is_confirmed", sa.Boolean(), nullable=True, server_default="0"),
        sa.Column("mfa_enabled", sa.Boolean(), nullable=True, server_default="0"),
        sa.Column("created_at", postgresql.TIMESTAMP(timezone=True), nullable=True),
        sa.Column("updated_at", postgresql.TIMESTAMP(timezone=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index("ix_users_full_name", "users", [sa.text("lower(full_name)")])
    op.create_index("ix_users_is_active", "users", ["is_active"])
    op.create_index("ix_users_email", "users", [sa.text("lower(email)")], unique=True)


def downgrade():
    op.drop_table("users")
