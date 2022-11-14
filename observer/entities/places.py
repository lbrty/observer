from pydantic import BaseModel

from observer.common.types import Identifier


class Country(BaseModel):
    id: Identifier
    name: str
    code: str


class State(BaseModel):
    id: Identifier
    name: str
    code: str
    country_id: Identifier


class City(BaseModel):
    id: Identifier
    name: str
    code: str
    country_id: Identifier
    state_id: Identifier
