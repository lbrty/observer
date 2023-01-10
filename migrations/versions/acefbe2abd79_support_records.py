"""support_records

Revision ID: acefbe2abd79
Revises: 98cf562f08e2
Create Date: 2022-10-09 22:08:36.117371
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "acefbe2abd79"
down_revision = "98cf562f08e2"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "support_records",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("description", sa.Text(), nullable=True),
        sa.Column("type", sa.Text(), nullable=False),
        sa.Column("consultant_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("age_group", sa.Text(), nullable=True),
        sa.Column("record_for", sa.Text(), nullable=False),
        sa.Column("owner_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("project_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("created_at", sa.DateTime(timezone=True), server_default=sa.text("now()"), nullable=True),
        sa.PrimaryKeyConstraint("id"),
        sa.ForeignKeyConstraint(
            ("consultant_id",),
            ["users.id"],
            ondelete="SET NULL",
        ),
        sa.ForeignKeyConstraint(
            ("project_id",),
            ["projects.id"],
            ondelete="SET NULL",
        ),
        sa.CheckConstraint("type IN ('humanitarian', 'legal', 'medical', 'general')", name="support_records_types"),
        sa.CheckConstraint("record_for IN ('person', 'pet')", name="support_records_record_for"),
        sa.CheckConstraint(
            """age_group IN (
                'infant',
                'toddler',
                'pre_school',
                'middle_childhood',
                'young_teen',
                'teenager',
                'young_adult',
                'early_adult',
                'middle_aged_adult',
                'old_adult'
            )""",
            name="support_records_age_group",
        ),
    )

    op.create_index(op.f("ix_support_records_type"), "support_records", ["type"])
    op.create_index(op.f("ix_support_records_description"), "support_records", [sa.text("lower(description)")])
    op.create_index(op.f("ix_support_records_consultant_id"), "support_records", ["consultant_id"])
    op.create_index(op.f("ix_support_records_beneficiary_age"), "support_records", [sa.text("beneficiary_age")])
    op.create_index(op.f("ix_support_records_owner_id"), "support_records", ["owner_id"])
    op.create_index(op.f("ix_support_records_project_id"), "support_records", ["project_id"])


def downgrade():
    op.drop_table("support_records")
