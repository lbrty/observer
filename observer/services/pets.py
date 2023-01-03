from typing import Optional, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.pets import NewPet, Pet, UpdatePet
from observer.repositories.pets import IPetsRepository
from observer.schemas.pets import NewPetRequest, UpdatePetRequest


class IPetsService(Protocol):
    repo: IPetsRepository

    async def create_pet(self, new_pet: NewPetRequest) -> Pet:
        raise NotImplementedError

    async def get_pet(self, pet_id: Identifier) -> Optional[Pet]:
        raise NotImplementedError

    async def update_pet(self, pet_id: Identifier, updates: UpdatePetRequest) -> Pet:
        raise NotImplementedError

    async def delete_pet(self, pet_id: Identifier) -> Pet:
        raise NotImplementedError


class PetsRepository(IPetsService):
    def __init__(self, repo: IPetsRepository):
        self.repo = repo

    async def create_pet(self, new_pet: NewPetRequest) -> Pet:
        pet = await self.repo.create_pet(NewPet(**new_pet.dict()))
        return pet

    async def get_pet(self, pet_id: Identifier) -> Pet:
        if pet := await self.repo.get_pet(pet_id):
            return pet

        raise NotFoundError(message="Pet not found")

    async def update_pet(self, pet_id: Identifier, updates: UpdatePetRequest) -> Pet:
        await self.get_pet(pet_id)
        pet = await self.repo.update_pet(pet_id, UpdatePet(**updates.dict()))
        return pet

    async def delete_pet(self, pet_id: Identifier) -> Pet:
        await self.get_pet(pet_id)
        pet = await self.repo.delete_pet(pet_id)
        return pet
