"""migration history

Revision ID: cdb2f0a93cea
Revises: d02d33dee75d
Create Date: 2022-12-25 20:38:27.400520
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "cdb2f0a93cea"
down_revision = "d02d33dee75d"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "migration_history",
        sa.Column(
            "id",
            postgresql.UUID(as_uuid=True),
            server_default=sa.text("gen_random_uuid()"),
            nullable=False,
        ),
        sa.Column("idp_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("migration_date", sa.DATE(), nullable=True),
        sa.Column("project_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("from_place_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("current_place_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column(
            "created_at",
            sa.DateTime(timezone=True),
            server_default=sa.text("now()"),
            nullable=True,
        ),
        sa.PrimaryKeyConstraint("id"),
        # If person deleted we no longer need to keep migration info
        sa.ForeignKeyConstraint(
            ("idp_id",),
            ["people.id"],
            ondelete="CASCADE",
        ),
        # If project gets deleted we still need to keep migration record
        sa.ForeignKeyConstraint(
            ("project_id",),
            ["projects.id"],
            ondelete="SET NULL",
        ),
        # If place was deleted we want to preserve record
        sa.ForeignKeyConstraint(
            ("from_place_id",),
            ["places.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("current_place_id",),
            ["places.id"],
            ondelete="SET NULL",
        ),
    )

    op.create_index(op.f("ix_migration_history_idp_id"), "migration_history", ["idp_id"])
    op.create_index(op.f("ix_migration_history_project_id"), "migration_history", ["project_id"])
    op.create_index(op.f("ix_migration_history_migration_date"), "migration_history", ["migration_date"])
    op.create_index(op.f("ix_migration_history_from_place_id"), "migration_history", ["from_place_id"])
    op.create_index(op.f("ix_migration_history_current_place_id"), "migration_history", ["current_place_id"])


def downgrade():
    op.drop_table("migration_history")
