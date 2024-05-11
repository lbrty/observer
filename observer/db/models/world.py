from sqlalchemy import ForeignKey, Text
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import mapped_column, relationship, Mapped, declared_attr

from observer.db.models import ModelBase


class BaseLocation(ModelBase):
    __abstract__ = True

    name: Mapped[str] = mapped_column(
        "name",
        Text(),
        nullable=False,
        index=True,
    )

    code: Mapped[str] = mapped_column(
        "code",
        Text(),
        nullable=False,
        unique=True,
    )


class Country(BaseLocation):
    __tablename__ = "countries"


class State(BaseLocation):
    __tablename__ = "states"
    country_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("countries.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    @declared_attr
    def country(self) -> Mapped[Country]:
        return relationship(
            Country,
            lazy="raise",
            back_populates="states",
            foreign_keys="[states.country_id]",
        )


class Place(BaseLocation):
    __tablename__ = "places"
    state_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("states.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    country_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("countries.id", ondelete="CASCADE"),
        nullable=False,
        index=True,
    )

    @declared_attr
    def state(self) -> Mapped[State]:
        return relationship(
            State,
            lazy="raise",
            back_populates="places",
            foreign_keys="[places.state_id]",
        )

    @declared_attr
    def country(self) -> Mapped[Country]:
        return relationship(
            Country,
            lazy="raise",
            back_populates="places",
            foreign_keys="[places.country_id]",
        )
