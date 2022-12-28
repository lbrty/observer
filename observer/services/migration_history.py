from typing import Protocol

from observer.common.types import Identifier
from observer.entities.migration_history import MigrationHistory
from observer.schemas.migration_history import (
    NewMigrationHistoryRequest,
    UpdateMigrationHistoryRequest,
)


class MigrationServiceInterface(Protocol):
    async def add_record(self, new_record: NewMigrationHistoryRequest) -> MigrationHistory:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError

    async def update_record(self, record_id: Identifier, updates: UpdateMigrationHistoryRequest) -> MigrationHistory:
        raise NotImplementedError

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError


class MigrationService(MigrationServiceInterface):
    async def add_record(self, new_record: NewMigrationHistoryRequest) -> MigrationHistory:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError

    async def update_record(self, record_id: Identifier, updates: UpdateMigrationHistoryRequest) -> MigrationHistory:
        raise NotImplementedError

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError
