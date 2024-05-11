from sqlalchemy import MetaData, UUID, text
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy.orm import Mapped, mapped_column

convention = {
    "ix": "ix_%(column_0_label)s",
    "uq": "uq_%(table_name)s_%(column_0_name)s_key",
    "ck": "ck_%(table_name)s_%(constraint_name)s",
    "fk": "fk_%(table_name)s_%(column_0_name)s_%(referred_table_name)s",
    "pk": "pk_%(table_name)s",
}

metadata = MetaData(naming_convention=convention)


class ModelBase(DeclarativeBase):
    __abstract__ = True
    metadata = metadata
    id: Mapped[UUID] = mapped_column(
        "id",
        UUID(as_uuid=True),
        primary_key=True,
        server_default=text("gen_random_uuid()"),
    )
