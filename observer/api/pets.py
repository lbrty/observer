from typing import List

from fastapi import APIRouter, Response, UploadFile
from starlette import status

from observer.schemas.documents import DocumentResponse
from observer.schemas.pets import (
    NewPetRequest,
    PetResponse,
    PetsResponse,
    UpdatePetRequest,
)

router = APIRouter(prefix="/pets")


@router.post(
    "",
    response_model=PetResponse,
    status_code=status.HTTP_201_CREATED,
    tags=["pets"],
)
async def create_pet(new_pet: NewPetRequest) -> PetResponse:
    ...


@router.get(
    "",
    response_model=PetsResponse,
    status_code=status.HTTP_200_OK,
    tags=["pets"],
)
async def get_pets() -> PetsResponse:
    ...


@router.get(
    "/{pet_id}",
    response_model=PetResponse,
    status_code=status.HTTP_200_OK,
    tags=["pets"],
)
async def get_pet() -> PetResponse:
    ...


@router.put(
    "/{pet_id}",
    response_model=PetResponse,
    status_code=status.HTTP_200_OK,
    tags=["pets"],
)
async def update_pet(updates: UpdatePetRequest) -> PetResponse:
    ...


@router.delete(
    "/{pet_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    tags=["pets"],
)
async def delete_pet() -> Response:
    ...


@router.post(
    "/{pet_id}/document",
    response_model=DocumentResponse,
    status_code=status.HTTP_201_CREATED,
    tags=["pets", "documents"],
)
async def pet_upload_document(file: UploadFile) -> DocumentResponse:
    ...


@router.get(
    "/{pet_id}/documents",
    response_model=List[DocumentResponse],
    status_code=status.HTTP_201_CREATED,
    tags=["pets", "documents"],
)
async def pet_get_documents() -> List[DocumentResponse]:
    ...
