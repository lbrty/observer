from typing import Optional, Protocol

from observer.common.types import Identifier
from observer.db import Database
from observer.entities.pets import NewPet, Pet, UpdatePet


class IPetsRepository(Protocol):
    async def create_pet(self, new_pet: NewPet) -> Pet:
        raise NotImplementedError

    async def get_pet(self, pet_id: Identifier) -> Optional[Pet]:
        raise NotImplementedError

    async def update_pet(self, pet_id: Identifier, updates: UpdatePet) -> Pet:
        raise NotImplementedError

    async def delete_pet(self, pet_id: Identifier) -> Optional[Pet]:
        raise NotImplementedError


class PetsRepository(IPetsRepository):
    def __init__(self, db: Database):
        self.db = db

    async def create_pet(self, new_pet: NewPet) -> Pet:
        raise NotImplementedError

    async def get_pet(self, pet_id: Identifier) -> Optional[Pet]:
        raise NotImplementedError

    async def update_pet(self, pet_id: Identifier, updates: UpdatePet) -> Pet:
        raise NotImplementedError

    async def delete_pet(self, pet_id: Identifier) -> Optional[Pet]:
        raise NotImplementedError
