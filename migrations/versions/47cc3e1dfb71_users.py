"""users

Revision ID: 47cc3e1dfb71
Revises: 268b7f793615
Create Date: 2022-10-08 21:13:19.224300
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "47cc3e1dfb71"
down_revision = "de0f408dabcf"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "users",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("email", sa.Text(), nullable=False),
        sa.Column("full_name", sa.Text(), nullable=True),
        sa.Column("password_hash", sa.Text(), nullable=False),
        sa.Column("role", sa.Text(), nullable=True),
        sa.Column("is_active", sa.Boolean(), nullable=True, server_default="1"),
        sa.Column("is_confirmed", sa.Boolean(), nullable=True, server_default="0"),
        sa.Column("office_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("mfa_enabled", sa.Boolean(), nullable=True, server_default="0"),
        sa.Column("mfa_encrypted_secret", sa.Text(), nullable=True),
        sa.Column("mfa_encrypted_backup_codes", sa.Text(), nullable=True),
        sa.PrimaryKeyConstraint("id"),
        sa.CheckConstraint("role IN ('admin', 'consultant', 'guest', 'staff')", name="users_role_type_check"),
        sa.ForeignKeyConstraint(
            ("office_id",),
            ["offices.id"],
            ondelete="SET NULL",
        ),
    )

    op.create_index(op.f("ix_users_full_name"), "users", [sa.text("lower(full_name)")])
    op.create_index(op.f("ix_users_is_active"), "users", ["is_active"])
    op.create_index(op.f("ix_users_office_id"), "users", ["office_id"]),
    op.create_index(op.f("ux_users_email"), "users", [sa.text("lower(email)")], unique=True)


def downgrade():
    op.drop_table("users")
