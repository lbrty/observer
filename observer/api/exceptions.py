from enum import Enum
from typing import Any

from starlette import status


class ErrorCode(str, Enum):
    unauthorized = "unauthorized"
    forbidden = "forbidden"
    not_found = "not_found"
    totp_error = "totp_error"
    totp_required_error = "totp_required_error"
    totp_exists_error = "totp_exists_error"
    registration_error = "registration_error"
    weak_password_error = "weak_password_error"
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


class InternalError(BaseAPIException):
    ...


class TOTPError(BaseAPIException):
    default_code = ErrorCode.totp_error
    default_status = status.HTTP_401_UNAUTHORIZED
    default_message = "invalid totp code"


class TOTPRequiredError(BaseAPIException):
    default_code = ErrorCode.totp_required_error
    default_status = status.HTTP_417_EXPECTATION_FAILED
    default_message = "totp required"


class TOTPExistsError(BaseAPIException):
    default_code = ErrorCode.totp_exists_error
    default_status = status.HTTP_400_BAD_REQUEST
    default_message = "totp already configured"


class UnauthorizedError(BaseAPIException):
    default_code = ErrorCode.unauthorized
    default_status = status.HTTP_401_UNAUTHORIZED
    default_message = "authentication required"


class ForbiddenError(BaseAPIException):
    default_code = ErrorCode.unauthorized
    default_status = status.HTTP_403_FORBIDDEN
    default_message = "access forbidden"


class NotFoundError(BaseAPIException):
    default_code = ErrorCode.not_found
    default_status = status.HTTP_404_NOT_FOUND
    default_message = "not found"


class RegistrationError(BaseAPIException):
    default_code = ErrorCode.registration_error
    default_status = status.HTTP_409_CONFLICT
    default_message = "registration error"


class WeakPasswordError(BaseAPIException):
    default_code = ErrorCode.weak_password_error
    default_status = status.HTTP_400_BAD_REQUEST
    default_message = "password is weak"
