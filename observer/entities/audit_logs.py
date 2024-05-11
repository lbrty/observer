from dataclasses import dataclass
from datetime import datetime
from typing import Any


@dataclass
class NewAuditLog:
    ref: str
    data: dict[str, Any]
    expires_at: datetime
