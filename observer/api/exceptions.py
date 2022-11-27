from enum import Enum
from typing import Any

from starlette import status


class ErrorCode(str, Enum):
    unauthorized = "unauthorized"
    forbidden = "forbidden"
    internal_error = "internal_error"


class BaseAPIException(Exception):
    default_code: ErrorCode = ErrorCode.internal_error
    default_message: str = "internal server error"
    default_status: int = status.HTTP_500_INTERNAL_SERVER_ERROR
    default_data: Any = None

    def __init__(self, code: ErrorCode = None, status_code: int = None, message: str = None, data: Any = None):
        self.code = code or self.default_code
        self.status = status_code or self.default_status
        self.message = message or self.default_message
        self.data = data or self.default_data

    def to_dict(self) -> dict:
        return dict(
            code=self.code.value,
            status_code=self.status,
            message=self.message,
            data=self.data,
        )


class UnauthorizedError(BaseAPIException):
    default_code = ErrorCode.unauthorized
    default_status = status.HTTP_401_UNAUTHORIZED


class ForbiddenError(BaseAPIException):
    default_code = ErrorCode.unauthorized
    default_status = status.HTTP_403_FORBIDDEN
