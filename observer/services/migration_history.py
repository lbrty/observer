from typing import List, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.migration_history import MigrationHistory, NewMigrationHistory
from observer.repositories.migration_history import IMigrationRepository
from observer.schemas.migration_history import NewMigrationHistoryRequest
from observer.services.world import IWorldService


class IMigrationService(Protocol):
    repo: IMigrationRepository
    world: IWorldService

    async def add_record(self, new_record: NewMigrationHistoryRequest) -> MigrationHistory:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError

    async def get_persons_records(self, person_id: Identifier) -> List[MigrationHistory]:
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

    async def get_persons_records(self, person_id: Identifier) -> List[MigrationHistory]:
        return await self.repo.get_persons_records(person_id)

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        if migration_history := await self.repo.delete_record(record_id):
            return migration_history

        raise NotFoundError(message="Migration record not found")
