from typing import Any, Dict, Union

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


def get_api_errors(*additional_responses) -> Dict[Union[int, str], Dict[str, Any]]:
    """Creates additional responses to swagger schema
    https://fastapi.tiangolo.com/advanced/additional-responses/?h=responses

    When called it expects `*additional_responses: tuple[int | str | tuple[Any, str]]` in the following format
    ```
    get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        (status.HTTP_404_NOT_FOUND, "Project or member not found"),
    )
    ```
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
