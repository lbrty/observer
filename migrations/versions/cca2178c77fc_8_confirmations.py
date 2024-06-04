"""confirmations

Revision ID: cca2178c77fc
Revises: 846f018af85c
Create Date: 2024-05-17 17:45:25.123875
"""

from typing import Sequence

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "cca2178c77fc"
down_revision: str | None = "846f018af85c"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.create_table(
        "confirmations",
        sa.Column(
            "id",
            sa.UUID(),
            server_default=sa.text("gen_random_uuid()"),
            nullable=False,
        ),
        sa.Column("code", sa.Text(), nullable=False),
        sa.Column("user_id", sa.UUID(), nullable=False),
        sa.Column("expires_at", sa.TIMESTAMP(timezone=True), nullable=True),
        sa.ForeignKeyConstraint(
            ["user_id"],
            ["users.id"],
            name=op.f("fk_confirmations_user_id_users"),
            ondelete="CASCADE",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_confirmations")),
        sa.UniqueConstraint("code", name=op.f("uq_confirmations_code_key")),
    )

    op.create_index(op.f("ix_confirmations_user_id"), "confirmations", ["user_id"])


def downgrade() -> None:
    op.drop_table("confirmations")