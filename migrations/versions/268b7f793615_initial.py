"""initial

Revision ID: 268b7f793615
Revises: 
Create Date: 2022-10-08 21:10:12.331589
"""
import sqlalchemy as sa
from alembic import op

# revision identifiers, used by Alembic.
revision = "268b7f793615"
down_revision = None
branch_labels = None
depends_on = None


def upgrade():
    op.execute(sa.text("CREATE EXTENSION IF NOT EXISTS pgcrypto"))


def downgrade():
    op.execute(sa.text("DROP EXTENSION IF EXISTS pgcrypto"))
