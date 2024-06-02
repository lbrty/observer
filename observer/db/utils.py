from datetime import datetime, UTC
from enum import Enum

from observer.common.reflect.inspect import unwrap_enum


def utcnow() -> datetime:
    return datetime.now(tz=UTC)


def choices_from_enum(enum: Enum) -> list[str]:
    values = list(unwrap_enum(enum))
    choices = [f"'{val}'" for val in values]
    return choices
