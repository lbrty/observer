"""initial

Revision ID: 40abefdabf70
Revises:
Create Date: 2024-05-10 17:28:44.599114
"""

from typing import Sequence

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "40abefdabf70"
down_revision: str | None = None
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.execute(sa.text("CREATE EXTENSION IF NOT EXISTS pgcrypto"))


def downgrade() -> None:
    op.execute(sa.text("DROP EXTENSION IF EXISTS pgcrypto"))
