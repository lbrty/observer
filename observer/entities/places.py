from pydantic import BaseModel

from observer.common.types import Identifier, PlaceType


class Country(BaseModel):
    id: Identifier
    name: str
    code: str


class State(BaseModel):
    id: Identifier
    name: str
    code: str
    country_id: Identifier


class Place(BaseModel):
    id: Identifier
    name: str
    code: str
    place_type: PlaceType
    country_id: Identifier
    state_id: Identifier
