from fastapi import Query

from observer.common.types import Pagination

DEFAULT_LIMIT = 100
MAX_LIMIT = 200
DEFAULT_OFFSET = 0


async def pagination(
    limit: int = Query(DEFAULT_LIMIT, description="How many items to return?"),
    offset: int = Query(DEFAULT_OFFSET, description="How many items to skip?"),
) -> Pagination:
    if limit < 0 or limit > MAX_LIMIT:
        limit = DEFAULT_LIMIT

    if offset < 0:
        offset = DEFAULT_OFFSET

    return Pagination(limit=limit, offset=offset)
