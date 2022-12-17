from pydantic import Field

from observer.schemas.base import SchemaBase


class Pagination(SchemaBase):
    limit: int = Field(100, description="How many items to show?")
    offset: int = Field(0, description="What is the starting point?")
