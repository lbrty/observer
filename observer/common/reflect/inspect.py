from enum import Enum
from typing import Any, Type, Dict


def unwrap_enum(val: Type[Enum]) -> Dict[str, Any]:
    """Get enum values as a dictionary"""
    return {v.name: v.value for v in val}
