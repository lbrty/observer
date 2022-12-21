from pydantic import BaseModel

from observer.common.types import Identifier, PlaceType


class BasePlace(BaseModel):
    name: str
    code: str


class Country(BasePlace):
    id: Identifier


class NewCountry(BasePlace):
    ...


class UpdateCountry(BasePlace):
    ...


class State(BasePlace):
    id: Identifier
    country_id: Identifier


class NewState(BasePlace):
    country_id: Identifier


class UpdateState(BasePlace):
    country_id: Identifier


class Place(BasePlace):
    id: Identifier
    place_type: PlaceType
    country_id: Identifier
    state_id: Identifier


class NewPlace(BasePlace):
    place_type: PlaceType
    country_id: Identifier
    state_id: Identifier


class UpdatePlace(BasePlace):
    place_type: PlaceType
    country_id: Identifier
    state_id: Identifier
