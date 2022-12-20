from datetime import datetime

from pydantic import BaseModel

from observer.common.types import Identifier, SomeStr


class Document(BaseModel):
    id: Identifier
    encryption_key: SomeStr
    name: str
    path: str
    owner_id: Identifier
    created_at: datetime
