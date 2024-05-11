from sqlalchemy import Text
from sqlalchemy.orm import Mapped, mapped_column

from observer.db.models import ModelBase


class Office(ModelBase):
    __tablename__ = "offices"

    name: Mapped[str] = mapped_column(
        "name",
        Text(),
        nullable=False,
        unique=True,
    )
