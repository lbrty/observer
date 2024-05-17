from sqlalchemy import Text, UUID, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column

from observer.db.models import ModelBase


class Project(ModelBase):
    __tablename__ = "projects"

    name: Mapped[str] = mapped_column(
        "name",
        Text(),
        nullable=False,
    )

    description: Mapped[str] = mapped_column(
        "description",
        Text(),
        nullable=True,
    )

    owner_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )
