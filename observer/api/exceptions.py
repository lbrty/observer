from enum import Enum, auto
from typing import Any

from starlette import status


class ErrorCode(str, Enum):
    unauthorized = "unauthorized"
    internal_error = "internal_error"


class BaseAPIException(BaseException):
    default_code: ErrorCode = ErrorCode.internal_error
    default_message: str = "internal server error"
    default_status: int = status.HTTP_500_INTERNAL_SERVER_ERROR
    default_data: Any = None

    def __init__(self, code: ErrorCode = None, status_code: int = None, message: str = None, data: Any = None):
        self.default_code = code or self.default_code
        self.default_status = status_code or self.default_status
        self.default_message = message or self.default_message
        self.data = data or self.default_data

    def to_dict(self) -> dict:
        return dict(
            code=self.default_code.value, status_code=self.default_status, message=self.default_message, data=self.data
        )


class UnauthorizedError(BaseAPIException):
    default_code = ErrorCode.unauthorized
    default_status = status.HTTP_401_UNAUTHORIZED
