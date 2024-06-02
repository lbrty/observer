from enum import Enum
from typing import Any


def unwrap_enum(val: Enum) -> dict[str, Any]:
    """Get enum values as a dictionary"""
    return {v.name: v.value for v in val}
