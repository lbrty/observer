from typing import Optional, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.support_records import (
    NewSupportRecord,
    SupportRecord,
    UpdateSupportRecord,
)
from observer.repositories.support_records import ISupportRecordsRepository
from observer.schemas.support_records import (
    NewSupportRecordRequest,
    UpdateSupportRecordRequest,
)


class ISupportRecordsService(Protocol):
    repo: ISupportRecordsRepository

    async def create_record(self, new_record: NewSupportRecordRequest) -> SupportRecord:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> Optional[SupportRecord]:
        raise NotImplementedError

    async def update_record(self, record_id: Identifier, updates: UpdateSupportRecordRequest) -> SupportRecord:
        raise NotImplementedError

    async def delete_record(self, record_id: Identifier) -> SupportRecord:
        raise NotImplementedError


class SupportRecordsService(ISupportRecordsService):
    def __init__(self, repo: ISupportRecordsRepository):
        self.repo = repo

    async def create_record(self, new_record: NewSupportRecordRequest) -> SupportRecord:
        return await self.repo.create_record(NewSupportRecord(**new_record.dict()))

    async def get_record(self, record_id: Identifier) -> Optional[SupportRecord]:
        if record := await self.repo.get_record(record_id):
            return record

        raise NotFoundError(message="Support record not found")

    async def update_record(self, record_id: Identifier, updates: UpdateSupportRecordRequest) -> SupportRecord:
        return await self.repo.update_record(record_id, UpdateSupportRecord(**updates.dict()))

    async def delete_record(self, record_id: Identifier) -> SupportRecord:
        return await self.repo.delete_record(record_id)
