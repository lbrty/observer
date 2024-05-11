"""offices

Revision ID: c685c35f23b2
Revises: 35830d41c820
Create Date: 2024-05-11 13:07:09.097278
"""

from typing import Sequence, Union

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "c685c35f23b2"
down_revision: Union[str, None] = "35830d41c820"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "offices",
        sa.Column(
            "id",
            sa.UUID(),
            server_default=sa.text("gen_random_uuid()"),
            nullable=False,
        ),
        sa.Column("name", sa.Text(), nullable=False),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_offices")),
        sa.UniqueConstraint("name", name=op.f("uq_offices_name_key")),
    )


def downgrade() -> None:
    op.drop_table("offices")
