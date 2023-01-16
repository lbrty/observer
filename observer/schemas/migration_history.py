from datetime import date, datetime
from typing import Optional

from pydantic import Field

from observer.common.types import Identifier
from observer.schemas.base import SchemaBase
from observer.schemas.world import PlaceResponse


class BaseMigrationHistory(SchemaBase):
    person_id: Identifier = Field(..., description="Person ID")
    migration_date: Optional[date] = Field(..., description="Date of migration")
    project_id: Identifier = Field(..., description="Project ID to which it belongs")
    from_place_id: Optional[Identifier] = Field(None, description="From place ID where then person has moved")
    current_place_id: Optional[Identifier] = Field(None, description="Current place of living")


class MigrationHistoryResponse(BaseMigrationHistory):
    id: Identifier = Field(..., description="Migration record ID")
    created_at: datetime = Field(..., description="Creation date")


class FullMigrationHistoryResponse(SchemaBase):
    id: Identifier = Field(..., description="Migration record ID")
    person_id: Identifier = Field(..., description="Person ID")
    migration_date: Optional[date] = Field(..., description="Date of migration")
    project_id: Identifier = Field(..., description="Project ID to which it belongs")
    from_place: Optional[PlaceResponse] = Field(None, description="Place where person has moved from")
    current_place: Optional[PlaceResponse] = Field(None, description="Current place of living")
    created_at: datetime = Field(..., description="Creation date")


class NewMigrationHistoryRequest(BaseMigrationHistory):
    ...


class UpdateMigrationHistoryRequest(BaseMigrationHistory):
    ...
