"""audit_logs

Revision ID: 0a7abb246c5a
Revises: 40abefdabf70
Create Date: 2024-05-11 12:21:30.124701
"""
from typing import Sequence, Union

import sqlalchemy as sa

from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision: str = "0a7abb246c5a"
down_revision: Union[str, None] = "40abefdabf70"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "audit_logs",
        sa.Column(
            "id", sa.UUID(), server_default=sa.text("gen_random_uuid()"), nullable=False
        ),
        sa.Column("ref", sa.Text(), nullable=False),
        sa.Column("data", postgresql.JSONB(astext_type=sa.Text()), nullable=True),
        sa.Column("created_at", sa.TIMESTAMP(timezone=True), nullable=False),
        sa.Column("expires_at", sa.TIMESTAMP(timezone=True), nullable=True),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_audit_logs")),
    )
    op.create_index(
        op.f("ix_audit_logs_expires_at"),
        "audit_logs",
        ["expires_at"],
        unique=False,
    )
    op.create_index(
        op.f("ix_audit_logs_ref"),
        "audit_logs",
        ["ref"],
        unique=False,
    )


def downgrade() -> None:
    op.drop_index(op.f("ix_audit_logs_ref"), table_name="audit_logs")
    op.drop_index(op.f("ix_audit_logs_expires_at"), table_name="audit_logs")
    op.drop_table("audit_logs")
