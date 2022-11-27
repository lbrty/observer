from fastapi.requests import Request
from fastapi.responses import JSONResponse

from observer.api.exceptions import BaseAPIException


async def handle_api_exception(_: Request, exc: BaseAPIException) -> JSONResponse:
    return JSONResponse(
        exc.to_dict(),
        status_code=exc.status,
    )
