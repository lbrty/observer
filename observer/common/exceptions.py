from typing import Any, Dict, Sequence, Tuple, TypeAlias

from fastapi.requests import Request
from fastapi.responses import JSONResponse
from pydantic.main import BaseModel

from observer.api.exceptions import BaseAPIException, ErrorCode

AdditionalResponses: TypeAlias = Sequence[int | Tuple[int, str]]


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


def get_api_errors(*additional_responses: AdditionalResponses) -> Dict[int, Dict]:
    """Creates additional responses to swagger schema

    https://fastapi.tiangolo.com/advanced/additional-responses/?h=responses
    """
    responses = {}
    for resp in additional_responses:
        response = dict(model=APIError)
        status_code = resp

        if isinstance(resp, tuple):
            status_code, description = resp
            response["description"] = description

        responses[status_code] = response

    return responses
