from uuid import UUID


def serialize_uuid_fields(data: dict) -> dict:
    return {field: str(value) if isinstance(value, UUID) else value for field, value in data.items()}
