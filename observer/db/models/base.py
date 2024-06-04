from sqlalchemy import MetaData, UUID, TIMESTAMP, text
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy.orm import Mapped, mapped_column

from observer.common.sqlalchemy import utc_now


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


class TimestampedModel(ModelBase):
    __abstract__ = True

    created_at: Mapped[TIMESTAMP] = mapped_column(
        "created_at",
        TIMESTAMP(timezone=True),
        default=utc_now,
    )

    updated_at: Mapped[TIMESTAMP] = mapped_column(
        "updated_at",
        TIMESTAMP(timezone=True),
        default=utc_now,
        onupdate=utc_now,
    )
