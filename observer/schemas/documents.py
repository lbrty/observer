from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, Field

from observer.common.types import Identifier
from observer.schemas.base import SchemaBase


class BaseDocument(SchemaBase):
    name: str = Field(..., description="Document filename")
    mimetype: str = Field(..., description="Document type (mimetype)")
    owner_id: Identifier = Field(..., description="Owner ID")
    project_id: Optional[Identifier] = Field(None, description="Project ID")


class DocumentResponse(BaseDocument):
    id: Identifier = Field(..., description="Document ID")
    created_at: datetime = Field(..., description="Creation date and time")


class DocumentsResponse(BaseModel):
    total: int = Field(..., description="Total count of documents")
    items: List[DocumentResponse] = Field(..., description="List of documents")


class NewDocumentRequest(BaseDocument):
    path: str = Field(..., description="Document path")
