"""password resets

Revision ID: 95fe34dc85fe
Revises: c6a0d931a680
Create Date: 2024-05-17 15:32:06.826509
"""

from typing import Sequence, Union

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "95fe34dc85fe"
down_revision: Union[str, None] = "c6a0d931a680"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "password_resets",
        sa.Column("code", sa.Text(), nullable=False),
        sa.Column("user_id", sa.UUID(), nullable=False),
        sa.Column("created_at", sa.TIMESTAMP(timezone=True), nullable=False),
        sa.Column(
            "id",
            sa.UUID(),
            server_default=sa.text("gen_random_uuid()"),
            nullable=False,
        ),
        sa.ForeignKeyConstraint(
            ["user_id"],
            ["users.id"],
            name=op.f("fk_password_resets_user_id_users"),
            ondelete="CASCADE",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_password_resets")),
    )
    op.create_index(
        op.f("ix_password_resets_code"),
        "password_resets",
        ["code"],
        unique=True,
    )
    op.create_index(
        op.f("ix_password_resets_user_id"),
        "password_resets",
        ["user_id"],
        unique=False,
    )


def downgrade() -> None:
    op.drop_index(op.f("ix_password_resets_user_id"), table_name="password_resets")
    op.drop_index(op.f("ix_password_resets_code"), table_name="password_resets")
    op.drop_table("password_resets")
