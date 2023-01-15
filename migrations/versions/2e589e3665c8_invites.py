"""invites

Revision ID: 2e589e3665c8
Revises: 06db81ceed60
Create Date: 2023-01-15 14:07:00.881716
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "2e589e3665c8"
down_revision = "06db81ceed60"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "invites",
        sa.Column("code", sa.Text()),
        sa.Column("user_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("expires_at", sa.DateTime(timezone=True), server_default=sa.text("now()"), nullable=False),
        sa.ForeignKeyConstraint(
            ("user_id",),
            ["users.id"],
            ondelete="CASCADE",
        ),
    )
    op.create_index(op.f("ux_invites_code"), "invites", ["code"], unique=True)
    op.create_index(op.f("ix_invites_user_id"), "invites", ["user_id"])


def downgrade():
    op.drop_table("invites")
