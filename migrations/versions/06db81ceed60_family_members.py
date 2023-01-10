"""family members

Revision ID: 06db81ceed60
Revises: cdb2f0a93cea
Create Date: 2023-01-10 20:39:52.317734
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = "06db81ceed60"
down_revision = "cdb2f0a93cea"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "family_members",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("age_group", sa.Text(), nullable=False),
        sa.Column("birth_date", sa.DATE(), nullable=True),
        sa.Column("sex", sa.Text(), nullable=True),
        sa.Column("notes", sa.Text(), nullable=True),
        sa.Column("idp_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("project_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.ForeignKeyConstraint(
            ("idp_id",),
            ["people.id"],
            ondelete="CASCADE",
        ),
        sa.ForeignKeyConstraint(
            ("project_id",),
            ["projects.id"],
            ondelete="SET NULL",
        ),
        sa.PrimaryKeyConstraint("id"),
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
            name="family_members_age_group",
        ),
        sa.CheckConstraint(
            """sex IN (
                'male',
                'female',
                'unknown'
            )""",
            name="family_members_sex",
        ),
    )

    op.create_index(op.f("family_members_full_name"), "family_members", [sa.text("lower(full_name)")])
    op.create_index(op.f("family_members_age_group"), "family_members", ["age_group"])
    op.create_index(op.f("family_members_birth_date"), "family_members", ["birth_date"])
    op.create_index(op.f("family_members_sex"), "family_members", ["sex"])
    op.create_index(op.f("family_members_idp_id"), "family_members", ["idp_id"])
    op.create_index(op.f("family_members_project_id"), "family_members", ["project_id"])


def downgrade():
    op.drop_table("family_members")
