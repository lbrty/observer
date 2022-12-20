from pydantic import BaseModel, Field

from observer.common.types import Identifier, PlaceType


# Countries
class BaseCountry(BaseModel):
    name: str = Field(..., description="Name of country")
    code: str = Field(..., description="Country code")


class CountryResponse(BaseCountry):
    id: Identifier = Field(..., description="Country ID")


class NewCountryRequest(BaseCountry):
    ...


class UpdateCountryRequest(BaseCountry):
    ...


# States
class BaseState(BaseModel):
    name: str = Field(..., description="Name of state")
    code: str = Field(..., description="State code")


class StateResponse(BaseState):
    id: Identifier = Field(..., description="State ID")
    country_id: Identifier = Field(..., description="Country ID")


class NewStateRequest(BaseState):
    country_id: Identifier = Field(..., description="Country ID")


class UpdateStateRequest(BaseCountry):
    country_id: Identifier = Field(..., description="Country ID")


# Places
class BasePlace(BaseModel):
    name: str = Field(..., description="Name of place")
    code: str = Field(..., description="Place code")
    place_type: PlaceType = Field(PlaceType.city, description="Type of place")


class PlaceResponse(BasePlace):
    id: Identifier = Field(..., description="Place ID")
    country_id: Identifier = Field(..., description="Country ID")
    state_id: Identifier = Field(..., description="State ID")


class NewPlaceRequest(BaseState):
    country_id: Identifier = Field(..., description="Country ID")
    state_id: Identifier = Field(..., description="State ID")


class UpdatePlaceRequest(BaseCountry):
    country_id: Identifier = Field(..., description="Country ID")
    state_id: Identifier = Field(..., description="State ID")
