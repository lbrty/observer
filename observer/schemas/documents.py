from datetime import datetime

from pydantic import BaseModel, Field

from observer.common.types import Identifier
from observer.schemas.base import SchemaBase


class BaseDocument(SchemaBase):
    encryption_key: str | None = Field(None, description="Encryption key")
    name: str = Field(..., description="Document filename")
    path: str = Field(..., description="Document path")
    mimetype: str = Field(..., description="Document type (mimetype)")
    owner_id: Identifier = Field(..., description="Owner ID")
    created_at: datetime = Field(..., description="Creation date and time")


class Document(BaseDocument):
    id: Identifier = Field(..., description="Document ID")


class DocumentsResponse(BaseModel):
    total: int = Field(..., description="Total count of documents")
    items: list[Document] = Field(..., description="List of documents")
