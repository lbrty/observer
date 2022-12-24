"""people

Revision ID: e07cd90eaa03
Revises: 193c427c31b7
Create Date: 2022-10-08 22:16:24.340875
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "e07cd90eaa03"
down_revision = "193c427c31b7"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "people",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
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
        sa.Column("from_place_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("current_place_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("project_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("category_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("consultant_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("tags", postgresql.ARRAY(sa.Text()), nullable=True),
        sa.Column("created_at", sa.DateTime(timezone=True), server_default=sa.text("now()"), nullable=True),
        sa.Column("updated_at", sa.DateTime(timezone=True), server_default=sa.text("now()"), nullable=True),
        sa.ForeignKeyConstraint(
            ("category_id",),
            ["categories.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("consultant_id",),
            ["users.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("current_place_id",),
            ["places.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("from_place_id",),
            ["places.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("project_id",),
            ["projects.id"],
            ondelete="SET NULL",
        ),
        sa.PrimaryKeyConstraint("id"),
        sa.CheckConstraint(
            """status IN (
                'consulted',
                'needs_call_back',
                'needs_legal_support',
                'needs_social_support',
                'needs_monitoring',
                'registered',
                'unknown'
            )""",
            name="people_status",
        ),
    )

    op.create_index(op.f("ix_people_full_name"), "people", [sa.text("lower(full_name)")])
    op.create_index(op.f("ix_people_reference_id"), "people", ["reference_id"])
    op.create_index(op.f("ix_people_status"), "people", ["status"])
    op.create_index(op.f("ix_people_email"), "people", [sa.text("lower(email)")])
    op.create_index(op.f("ix_people_birth_date"), "people", ["birth_date"])
    op.create_index(op.f("ix_people_category_id"), "people", ["category_id"])
    op.create_index(op.f("ix_people_consultant_id"), "people", ["consultant_id"])
    op.create_index(op.f("is_people_current_place_id"), "people", ["current_place_id"])
    op.create_index(op.f("ix_people_from_place_id"), "people", ["from_place_id"])
    op.create_index(op.f("ix_people_project_id"), "people", ["project_id"])
    op.create_index(op.f("ix_people_tags"), "people", ["tags"])


def downgrade():
    op.drop_table("people")
