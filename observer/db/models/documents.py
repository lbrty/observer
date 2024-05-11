# from sqlalchemy import Text, Integer, UUID, ForeignKey
# from sqlalchemy.orm import Mapped, mapped_column
#
# from observer.db.models import ModelBase
#
#
# class Category(ModelBase):
#     __tablename__ = "documents"
#
#     encryption_key: Mapped[str | None] = mapped_column(
#         "encryption_key",
#         Text(),
#         nullable=False,
#     )
#
#     name: Mapped[str] = mapped_column(
#         "name",
#         Text(),
#         nullable=False,
#         index=True,
#     )
#
#     path: Mapped[str] = mapped_column(
#         "path",
#         Text(),
#         nullable=False,
#     )
#
#     mime: Mapped[str] = mapped_column(
#         "mime",
#         Text(),
#         nullable=False,
#     )
#
#     size: Mapped[float] = mapped_column(
#         "size",
#         Integer(),
#         nullable=False,
#     )
#
#     owner_id: Mapped[UUID] = mapped_column(
#         "owner_id",
#         UUID(as_uuid=True),
#         nullable=False,
#         index=True
#     )
#     owner: Mapped[User] = mapped_column(
#         "owner",
#         ForeignKey("users.id", ondelete="CASCADE"),
#         nullable=False,
#     )
#     project_id: Mapped[UUID] = mapped_column(
#         "project_id",
#         UUID(as_uuid=True),
#         nullable=False,
#         index=True
#     )
#
#     project: Mapped[Project] = mapped_column(
#         "project",
#         ForeignKey("projects.id", ondelete="CASCADE"),
#         nullable=False,
#     )
