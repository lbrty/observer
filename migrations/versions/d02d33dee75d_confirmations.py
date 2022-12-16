"""confirmations

Revision ID: d02d33dee75d
Revises: fb2c4d5a1f31
Create Date: 2022-12-12 13:25:20.514063
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "d02d33dee75d"
down_revision = "fb2c4d5a1f31"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "confirmations",
        sa.Column("code", sa.Text()),
        sa.Column("user_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("expires_at", sa.DateTime(timezone=True), server_default=sa.text("now()"), nullable=True),
        sa.ForeignKeyConstraint(
            ("user_id",),
            ["users.id"],
            ondelete="CASCADE",
        ),
    )
    op.create_index(op.f("ux_confirmations_code"), "confirmations", ["code"], unique=True)
    op.create_index(op.f("ix_confirmations_user_id"), "confirmations", ["user_id"])


def downgrade():
    op.drop_table("confirmations")
