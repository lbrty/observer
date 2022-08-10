"""initial extensions

Revision ID: 65bd3ba22521
Revises:
Create Date: 2022-08-07 22:21:29.485044
"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = "65bd3ba22521"
down_revision = None
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.execute(sa.text("CREATE EXTENSION IF NOT EXISTS pgcrypto"))


def downgrade() -> None:
    op.execute(sa.text("DROP EXTENSION IF EXISTS pgcrypto"))
