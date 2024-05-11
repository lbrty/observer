"""categories

Revision ID: 35830d41c820
Revises: 0a7abb246c5a
Create Date: 2024-05-11 12:41:01.983388
"""

from typing import Sequence, Union

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "35830d41c820"
down_revision: Union[str, None] = "0a7abb246c5a"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "categories",
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column(
            "id", sa.UUID(), server_default=sa.text("gen_random_uuid()"), nullable=False
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_categories")),
        sa.UniqueConstraint("name", name=op.f("uq_categories_name_key")),
    )


def downgrade() -> None:
    op.drop_table("categories")
