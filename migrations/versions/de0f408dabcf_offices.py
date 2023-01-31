"""offices

Revision ID: de0f408dabcf
Revises: 2e589e3665c8
Create Date: 2022-10-08 21:10:10.444300
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "de0f408dabcf"
down_revision = "268b7f793615"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "offices",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("name", sa.Text(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )


def downgrade():
    op.drop_table("offices")
