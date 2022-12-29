from typing import List, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.migration_history import (
    MigrationHistory,
    NewMigrationHistory,
    UpdateMigrationHistory,
)
from observer.repositories.migration_history import IMigrationRepository
from observer.schemas.migration_history import (
    NewMigrationHistoryRequest,
    UpdateMigrationHistoryRequest,
)
from observer.services.world import IWorldService


class IMigrationService(Protocol):
    repo: IMigrationRepository
    world: IWorldService

    async def add_record(self, new_record: NewMigrationHistoryRequest) -> MigrationHistory:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError

    async def get_records(self, record_id: Identifier) -> List[MigrationHistory]:
        raise NotImplementedError

    async def update_record(self, record_id: Identifier, updates: UpdateMigrationHistoryRequest) -> MigrationHistory:
        raise NotImplementedError

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError


class MigrationService(IMigrationService):
    def __init__(self, repo: IMigrationRepository, world: IWorldService):
        self.repo = repo
        self.world = world

    async def add_record(self, new_record: NewMigrationHistoryRequest) -> MigrationHistory:
        if new_record.from_place_id:
            await self.world.get_place(new_record.from_place_id)

        if new_record.current_place_id:
            await self.world.get_place(new_record.current_place_id)

        return await self.repo.add_record(NewMigrationHistory(**new_record.dict()))

    async def get_record(self, record_id: Identifier) -> MigrationHistory:
        if migration_record := await self.repo.get_record(record_id):
            return migration_record

        raise NotFoundError(message="Migration record not found")

    async def get_records(self, record_id: Identifier) -> List[MigrationHistory]:
        return await self.repo.get_records(record_id)

    async def update_record(self, record_id: Identifier, updates: UpdateMigrationHistoryRequest) -> MigrationHistory:
        if updates.from_place_id:
            await self.world.get_place(updates.from_place_id)

        if updates.current_place_id:
            await self.world.get_place(updates.current_place_id)

        return await self.repo.update_record(record_id, UpdateMigrationHistory(**updates.dict()))

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        return await self.repo.delete_record(record_id)
