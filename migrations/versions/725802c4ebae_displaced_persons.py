"""displaced persons

Revision ID: 725802c4ebae
Revises: b4849b0f8d04
Create Date: 2022-08-10 19:59:40.736187
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "725802c4ebae"
down_revision = "b4849b0f8d04"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "displaced_persons",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("encryption_key", sa.Text(), nullable=True),
        sa.Column("status", sa.Text(), nullable=True),
        sa.Column("external_id", sa.Text(), nullable=True),
        sa.Column("reference_id", sa.Text(), nullable=True),
        sa.Column("email", sa.Text(), nullable=True),
        sa.Column("full_name", sa.Text(), nullable=False),
        sa.Column("birth_date", sa.DATE(), nullable=True),
        sa.Column("notes", sa.Text(), nullable=True),
        sa.Column("phone_number", sa.Text(), nullable=True),
        sa.Column("phone_number_additional", sa.Text(), nullable=True),
        sa.Column("migration_date", sa.DATE(), nullable=True),
        sa.Column("from_city_id", postgresql.UUID(), nullable=True),
        sa.Column("from_state_id", postgresql.UUID(), nullable=True),
        sa.Column("current_city_id", postgresql.UUID(), nullable=True),
        sa.Column("current_state_id", postgresql.UUID(), nullable=True),
        sa.Column("project_id", postgresql.UUID(), nullable=True),
        sa.Column("category_id", postgresql.UUID(), nullable=True),
        sa.Column("creator_id", postgresql.UUID(), nullable=True),
        sa.Column("tags", postgresql.ARRAY(sa.Text()), nullable=True),
        sa.Column("created_at", postgresql.TIMESTAMP(timezone=True)),
        sa.Column("updated_at", postgresql.TIMESTAMP(timezone=True)),
        sa.ForeignKeyConstraint(
            ("category_id",),
            ["categories.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("creator_id",),
            ["users.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("current_city_id",),
            ["cities.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("current_state_id",),
            ["states.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("from_city_id",),
            ["cities.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("from_state_id",),
            ["states.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("project_id",),
            ["projects.id"],
            ondelete="SET NULL",
        ),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ix_displaced_persons_full_name"), "displaced_persons", [sa.text("lower(full_name)")])
    op.create_index(op.f("ix_displaced_persons_status"), "displaced_persons", ["status"])
    op.create_index(op.f("ix_displaced_persons_email"), "displaced_persons", [sa.text("lower(email)")])
    op.create_index(op.f("ix_displaced_persons_birth_date"), "displaced_persons", ["birth_date"])
    op.create_index(op.f("ix_displaced_persons_tags"), "displaced_persons", ["tags"])


def downgrade():
    op.drop_table("displaced_persons")
