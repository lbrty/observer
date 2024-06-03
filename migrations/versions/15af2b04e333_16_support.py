"""support

Revision ID: 15af2b04e333
Revises: 06a42be25f0d
Create Date: 2024-06-03 23:13:39.892193
"""

from typing import Sequence

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "15af2b04e333"
down_revision: str | None = "06a42be25f0d"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.create_table(
        "support_records",
        sa.Column("id", sa.UUID(), server_default=sa.text("gen_random_uuid()"), nullable=False),
        sa.Column("notes", sa.Text(), nullable=True),
        sa.Column("type", sa.Text(), nullable=False),
        sa.Column("owner_id", sa.Text(), nullable=False),
        sa.Column("migration_date", sa.TIMESTAMP(timezone=True), nullable=True),
        sa.Column("consultant_id", sa.UUID(), nullable=True),
        sa.Column("person_id", sa.UUID(), nullable=False),
        sa.Column("project_id", sa.UUID(), nullable=True),
        sa.Column("created_at", sa.TIMESTAMP(timezone=True), nullable=False),
        sa.Column("updated_at", sa.TIMESTAMP(timezone=True), nullable=False),
        sa.CheckConstraint(
            "type IN ('humanitarian', 'legal', 'medical', 'general')",
            name=op.f("ck_support_records_type"),
        ),
        sa.ForeignKeyConstraint(
            ["consultant_id"],
            ["users.id"],
            name=op.f("fk_support_records_consultant_id_users"),
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ["person_id"],
            ["people.id"],
            name=op.f("fk_support_records_person_id_people"),
            ondelete="CASCADE",
        ),
        sa.ForeignKeyConstraint(
            ["project_id"],
            ["projects.id"],
            name=op.f("fk_support_records_project_id_projects"),
            ondelete="SET NULL",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_support_records")),
    )

    op.create_index(op.f("ix_support_records_owner_id"), "support_records", ["owner_id"])
    op.create_index(op.f("ix_support_records_consultant_id"), "support_records", ["consultant_id"])
    op.create_index(op.f("ix_support_records_person_id"), "support_records", ["person_id"])
    op.create_index(op.f("ix_support_records_project_id"), "support_records", ["project_id"])


def downgrade() -> None:
    op.drop_table("support_records")
