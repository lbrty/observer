from pydantic import BaseModel

from observer.common.types import Identifier


class Project(BaseModel):
    id: Identifier
    name: str
    description: str | None
