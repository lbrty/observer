from typing import Protocol

from observer.common.types import Identifier
from observer.entities.migration_history import (
    MigrationHistory,
    NewMigrationHistory,
    UpdateMigrationHistory,
)


class MigrationRepositoryInterface(Protocol):
    async def add_record(self, new_record: NewMigrationHistory) -> MigrationHistory:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError

    async def update_record(self, record_id: Identifier, updates: UpdateMigrationHistory) -> MigrationHistory:
        raise NotImplementedError

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError


class MigrationRepository(MigrationRepositoryInterface):
    async def add_record(self, new_record: NewMigrationHistory) -> MigrationHistory:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError

    async def update_record(self, record_id: Identifier, updates: UpdateMigrationHistory) -> MigrationHistory:
        raise NotImplementedError

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError
