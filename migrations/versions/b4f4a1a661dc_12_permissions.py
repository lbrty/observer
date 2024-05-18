"""permissions

Revision ID: b4f4a1a661dc
Revises: bcf9e218d4a0
Create Date: 2024-05-18 14:02:38.534731
"""

from typing import Sequence, Union

import sqlalchemy as sa

from alembic import op


# revision identifiers, used by Alembic.
revision: str = "b4f4a1a661dc"
down_revision: Union[str, None] = "bcf9e218d4a0"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "permissions",
        sa.Column(
            "id", sa.UUID(), server_default=sa.text("gen_random_uuid()"), nullable=False
        ),
        sa.Column("notes", sa.Text(), nullable=False),
        sa.Column("can_create", sa.Boolean(), nullable=False),
        sa.Column("can_read", sa.Boolean(), nullable=False),
        sa.Column("can_update", sa.Boolean(), nullable=False),
        sa.Column("can_delete", sa.Boolean(), nullable=False),
        sa.Column("can_read_documents", sa.Boolean(), nullable=False),
        sa.Column("can_read_personal_info", sa.Boolean(), nullable=False),
        sa.Column("can_invite_members", sa.Boolean(), nullable=False),
        sa.Column("user_id", sa.UUID(), nullable=False),
        sa.Column("project_id", sa.UUID(), nullable=False),
        sa.ForeignKeyConstraint(
            ["project_id"],
            ["projects.id"],
            name=op.f("fk_permissions_project_id_projects"),
            ondelete="CASCADE",
        ),
        sa.ForeignKeyConstraint(
            ["user_id"],
            ["users.id"],
            name=op.f("fk_permissions_user_id_users"),
            ondelete="CASCADE",
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_permissions")),
    )
    op.create_index(op.f("ix_permissions_project_id"), "permissions", ["project_id"])
    op.create_index(op.f("ix_permissions_user_id"), "permissions", ["user_id"])


def downgrade() -> None:
    op.drop_index(op.f("ix_permissions_user_id"), table_name="permissions")
    op.drop_index(op.f("ix_permissions_project_id"), table_name="permissions")
    op.drop_table("permissions")
