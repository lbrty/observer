from typing import List

from fastapi import APIRouter

from observer.schemas.places import CountryResponse, PlaceResponse, StateResponse

router = APIRouter(prefix="/places")


@router.get("/countries", response_model=List[CountryResponse])
async def get_countries() -> List[CountryResponse]:
    pass


@router.get("/states", response_model=List[StateResponse])
async def get_states() -> List[StateResponse]:
    pass


@router.get("/places", response_model=List[PlaceResponse])
async def get_states() -> List[PlaceResponse]:
    pass
