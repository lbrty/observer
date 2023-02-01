from pydantic import BaseModel

from observer.common.types import Identifier


class Office(BaseModel):
    id: Identifier
    name: str
