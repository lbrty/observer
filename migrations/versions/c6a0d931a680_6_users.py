"""users

Revision ID: c6a0d931a680
Revises: f5af5093abf5
Create Date: 2024-05-12 21:34:28.441090
"""

from typing import Sequence, Union

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "c6a0d931a680"
down_revision: Union[str, None] = "f5af5093abf5"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "users",
        sa.Column(
            "id", sa.UUID(), server_default=sa.text("gen_random_uuid()"), nullable=False
        ),
        sa.Column("email", sa.Text(), nullable=False),
        sa.Column("full_name", sa.Text(), nullable=True),
        sa.Column("password_hash", sa.Text(), nullable=False),
        sa.Column("role", sa.Text(), nullable=True),
        sa.Column("is_active", sa.Boolean(), nullable=True),
        sa.Column("is_confirmed", sa.Boolean(), nullable=True),
        sa.Column("office_id", sa.UUID(), nullable=True),
        sa.Column("mfa_enabled", sa.Boolean(), nullable=True),
        sa.Column("mfa_encrypted_secret", sa.Text(), nullable=True),
        sa.Column("mfa_encrypted_backup_codes", sa.Text(), nullable=True),
        sa.CheckConstraint(
            "role IN ('admin', 'staff', 'consultant', 'guest')",
            name=op.f("ck_users_users_role_type_check"),
        ),
        sa.ForeignKeyConstraint(
            ["office_id"],
            ["offices.id"],
            name=op.f("fk_users_office_id_offices"),
            ondelete="SET NULL",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_users")),
        sa.UniqueConstraint("email", name=op.f("uq_users_email_key")),
    )
    op.create_index(op.f("ix_users_office_id"), "users", ["office_id"])


def downgrade() -> None:
    op.drop_index(op.f("ix_users_office_id"), table_name="users")
    op.drop_table("users")
