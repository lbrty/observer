from typing import Any, Dict, List

from fastapi.requests import Request
from fastapi.responses import JSONResponse
from pydantic.main import BaseModel

from observer.api.exceptions import BaseAPIException, ErrorCode


class APIError(BaseModel):
    code: ErrorCode
    status_code: int
    message: str
    data: Any


async def handle_api_exception(_: Request, exc: BaseAPIException) -> JSONResponse:
    return JSONResponse(
        exc.to_dict(),
        status_code=exc.status,
    )


def get_api_errors(*status_codes: List[int]) -> Dict[int, Dict[str, APIError]]:
    return {status_code: dict(model=APIError) for status_code in status_codes}
