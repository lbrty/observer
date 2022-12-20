from pydantic import BaseModel, Field

from observer.common.types import Identifier, PlaceType


class CountryResponse(BaseModel):
    id: Identifier = Field(..., description="Country ID")
    name: str = Field(..., description="Name of country")
    code: str = Field(..., description="Country code")


class StateResponse(BaseModel):
    id: Identifier = Field(..., description="State ID")
    name: str = Field(..., description="Name of state")
    code: str = Field(..., description="State code")
    country_id: Identifier


class PlaceResponse(BaseModel):
    id: Identifier = Field(..., description="City ID")
    name: str = Field(..., description="Name of city")
    code: str = Field(..., description="City code")
    place_type: PlaceType = Field(PlaceType.city, description="Type of place")
    country_id: Identifier
    state_id: Identifier
