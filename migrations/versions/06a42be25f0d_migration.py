"""migration

Revision ID: 06a42be25f0d
Revises: 6f50b18ba1b7
Create Date: 2024-06-03 22:52:36.483276
"""

from typing import Sequence

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "06a42be25f0d"
down_revision: str | None = "6f50b18ba1b7"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.create_table(
        "migration_records",
        sa.Column("id", sa.UUID(), server_default=sa.text("gen_random_uuid()"), nullable=False),
        sa.Column("notes", sa.Text(), nullable=True),
        sa.Column("migration_date", sa.TIMESTAMP(timezone=True), nullable=True),
        sa.Column("person_id", sa.UUID(), nullable=False),
        sa.Column("project_id", sa.UUID(), nullable=False),
        sa.Column("current_place_id", sa.UUID(), nullable=False),
        sa.Column("from_place_id", sa.UUID(), nullable=False),
        sa.Column("created_at", sa.TIMESTAMP(timezone=True), nullable=False),
        sa.Column("updated_at", sa.TIMESTAMP(timezone=True), nullable=False),
        sa.ForeignKeyConstraint(
            ["current_place_id"],
            ["places.id"],
            name=op.f("fk_migration_records_current_place_id_places"),
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ["from_place_id"],
            ["places.id"],
            name=op.f("fk_migration_records_from_place_id_places"),
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ["person_id"],
            ["people.id"],
            name=op.f("fk_migration_records_person_id_people"),
            ondelete="CASCADE",
        ),
        sa.ForeignKeyConstraint(
            ["project_id"],
            ["projects.id"],
            name=op.f("fk_migration_records_project_id_projects"),
            ondelete="SET NULL",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_migration_records")),
    )

    op.create_index(op.f("ix_migration_records_person_id"), "migration_records", ["person_id"])
    op.create_index(op.f("ix_migration_records_project_id"), "migration_records", ["project_id"])
    op.create_index(
        op.f("ix_migration_records_current_place_id"), "migration_records", ["current_place_id"]
    )
    op.create_index(
        op.f("ix_migration_records_from_place_id"), "migration_records", ["from_place_id"]
    )


def downgrade() -> None:
    op.drop_table("migration_records")
