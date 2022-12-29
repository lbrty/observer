from datetime import date, datetime
from typing import Optional

from pydantic import Field

from observer.common.types import Identifier
from observer.schemas.base import SchemaBase


class BaseMigrationHistory(SchemaBase):
    idp_id: Identifier = Field(..., description="IDP ID")
    migration_date: Optional[date] = Field(..., description="Date of migration")
    project_id: Identifier = Field(..., description="Project ID to which it belongs")
    from_place_id: Identifier = Field(..., description="From place ID where IDP has moved")
    current_place_id: Identifier = Field(..., description="Current place of living")


class MigrationHistoryResponse(SchemaBase):
    id: Identifier = Field(..., description="Migration record ID")
    created_at: datetime = Field(..., description="Creation date")


class NewMigrationHistoryRequest(BaseMigrationHistory):
    ...


class UpdateMigrationHistoryRequest(BaseMigrationHistory):
    ...
