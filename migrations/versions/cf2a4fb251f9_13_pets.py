"""pets

Revision ID: cf2a4fb251f9
Revises: b4f4a1a661dc
Create Date: 2024-05-24 21:13:32.804526
"""

from typing import Sequence

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "cf2a4fb251f9"
down_revision: str | None = "b4f4a1a661dc"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.create_table(
        "pets",
        sa.Column("name", sa.Text(), nullable=False),
        sa.Column("notes", sa.Text(), nullable=True),
        sa.Column("status", sa.Text(), nullable=False),
        sa.Column("registration_id", sa.Text(), nullable=False),
        sa.Column("created_at", sa.TIMESTAMP(timezone=True), nullable=False),
        sa.Column("owner_id", sa.UUID(), nullable=False),
        sa.Column("project_id", sa.UUID(), nullable=False),
        sa.Column(
            "id",
            sa.UUID(),
            server_default=sa.text("gen_random_uuid()"),
            nullable=False,
        ),
        sa.CheckConstraint(
            "status IN ('registered', 'adopted', 'owner_found', 'needs_shelter', 'unknown')",
            name=op.f("ck_pets_status"),
        ),
        sa.ForeignKeyConstraint(
            ["owner_id"],
            ["users.id"],
            name=op.f("fk_pets_owner_id_users"),
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ["project_id"],
            ["projects.id"],
            name=op.f("fk_pets_project_id_projects"),
            ondelete="SET NULL",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_pets")),
    )

    op.create_index("ix_name_pets", "pets", [sa.text("lower(name)")])
    op.create_index(op.f("ix_pets_status"), "pets", ["status"])
    op.create_index(op.f("ix_pets_registration_id"), "pets", ["registration_id"])
    op.create_index(op.f("ix_pets_owner_id"), "pets", ["owner_id"])
    op.create_index(op.f("ix_pets_project_id"), "pets", ["project_id"])


def downgrade() -> None:
    op.drop_table("pets")
