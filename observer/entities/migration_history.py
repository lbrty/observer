from datetime import date, datetime
from typing import Optional

from observer.common.types import Identifier
from observer.entities.base import ModelBase


class BaseMigrationHistory(ModelBase):
    idp_id: Identifier
    migration_date: Optional[date]
    from_place_id: Optional[Identifier]
    current_place_id: Optional[Identifier]


class MigrationHistory(ModelBase):
    id: Identifier
    created_at: datetime


class NewMigrationHistory(BaseMigrationHistory):
    ...


class UpdateMigrationHistory(BaseMigrationHistory):
    ...
