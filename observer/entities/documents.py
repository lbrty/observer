from datetime import datetime
from typing import Dict

from pydantic import BaseModel

from observer.common.types import Identifier, SomeStr

AllowedDocumentTypes: Dict[str, str] = {
    # Images
    "image/jpeg": ".jpg",
    "image/png": ".png",
    # Documents
    "text/csv": ".csv",
    "text/plain": ".txt",
    "application/msword": ".doc",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
    "application/vnd.ms-excel": "xls",
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": "xlsx",
    "application/rtf": ".rtf",
    "application/pdf": ".pdf",
    # Media files
    "video/mp4": ".mp4",
    "video/mpeg": ".mpeg",
    "audio/mpeg": ".mp3",
}


class Document(BaseModel):
    id: Identifier
    encryption_key: SomeStr
    name: str
    path: str
    mimetype: str
    owner_id: Identifier
    created_at: datetime
