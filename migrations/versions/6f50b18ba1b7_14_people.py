"""people

Revision ID: 6f50b18ba1b7
Revises: cf2a4fb251f9
Create Date: 2024-06-02 17:42:03.416726
"""

from typing import Sequence

import sqlalchemy as sa

from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision: str = "6f50b18ba1b7"
down_revision: str | None = "cf2a4fb251f9"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None

age_groups = (
    "age_group IN ("
    "'infant', 'toddler', 'pre_school', 'middle_childhood', "
    "'young_teen', 'teenager', 'young_adult', 'early_adult', "
    "'middle_aged_adult', 'old_adult', 'unknown'"
    ")"
)

status = (
    "status IN ("
    "'consulted', 'needs_call_back', 'needs_legal_support', "
    "'needs_social_support', 'needs_monitoring', 'registered', 'unknown'"
    ")"
)


def upgrade() -> None:
    op.create_table(
        "people",
        sa.Column("id", sa.UUID(), server_default=sa.text("gen_random_uuid()"), nullable=False),
        sa.Column("status", sa.Text(), nullable=False),
        sa.Column("external_id", sa.Text(), nullable=True),
        sa.Column("email", sa.Text(), nullable=True),
        sa.Column("full_name", sa.Text(), nullable=True),
        sa.Column("birth_date", sa.Date(), nullable=True),
        sa.Column("sex", sa.Text(), nullable=True),
        sa.Column("notes", sa.Text(), nullable=True),
        sa.Column("phone_numbers", postgresql.JSONB(astext_type=sa.Text()), nullable=True),
        sa.Column("age_group", sa.Text(), nullable=False),
        sa.Column("project_id", sa.UUID(), nullable=True),
        sa.Column("category_id", sa.UUID(), nullable=True),
        sa.Column("parent_id", sa.UUID(), nullable=True),
        sa.Column("consultant_id", sa.UUID(), nullable=True),
        sa.Column("office_id", sa.UUID(), nullable=True),
        sa.Column("tags", postgresql.ARRAY(sa.Text()), nullable=True),
        sa.Column("created_at", sa.TIMESTAMP(timezone=True), nullable=False),
        sa.Column("updated_at", sa.TIMESTAMP(timezone=True), nullable=False),
        sa.CheckConstraint(age_groups, name=op.f("ck_people_age_group")),
        sa.CheckConstraint("sex IN ('male', 'female', 'unknown')", name=op.f("ck_people_sex")),
        sa.CheckConstraint(status, name=op.f("ck_people_status")),
        sa.ForeignKeyConstraint(
            ["category_id"],
            ["categories.id"],
            name=op.f("fk_people_category_id_categories"),
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ["consultant_id"],
            ["users.id"],
            name=op.f("fk_people_consultant_id_users"),
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ["office_id"],
            ["offices.id"],
            name=op.f("fk_people_office_id_offices"),
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ["parent_id"],
            ["people.id"],
            name=op.f("fk_people_parent_id_people"),
            ondelete="CASCADE",
        ),
        sa.ForeignKeyConstraint(
            ["project_id"],
            ["projects.id"],
            name=op.f("fk_people_project_id_projects"),
            ondelete="SET NULL",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_people")),
    )

    op.create_index("ix_people_email", "people", [sa.text("lower(email)")])
    op.create_index("ix_people_full_name", "people", [sa.text("lower(full_name)")])
    op.create_index(op.f("ix_people_project_id"), "people", ["project_id"])
    op.create_index(op.f("ix_people_category_id"), "people", ["category_id"])
    op.create_index(op.f("ix_people_parent_id"), "people", ["parent_id"])
    op.create_index(op.f("ix_people_consultant_id"), "people", ["consultant_id"])
    op.create_index(op.f("ix_people_office_id"), "people", ["office_id"])
    op.create_index(op.f("ix_people_tags"), "people", ["tags"])


def downgrade() -> None:
    op.drop_table("people")
