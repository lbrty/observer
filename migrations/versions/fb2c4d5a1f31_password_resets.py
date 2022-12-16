"""password resets

Revision ID: fb2c4d5a1f31
Revises: acefbe2abd79
Create Date: 2022-12-09 23:44:25.208805
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "fb2c4d5a1f31"
down_revision = "acefbe2abd79"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "password_resets",
        sa.Column("code", sa.Text()),
        sa.Column("user_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("created_at", sa.DateTime(timezone=True), server_default=sa.text("now()"), nullable=True),
        sa.ForeignKeyConstraint(
            ("user_id",),
            ["users.id"],
            ondelete="CASCADE",
        ),
    )
    op.create_index(op.f("ux_password_resets_code"), "password_resets", ["code"], unique=True)
    op.create_index(op.f("ix_password_resets_user_id"), "password_resets", ["user_id"])


def downgrade():
    op.drop_table("password_resets")
