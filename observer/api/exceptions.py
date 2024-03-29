from enum import Enum
from typing import Any, Dict, Optional

from starlette import status


class ErrorCode(str, Enum):
    unauthorized = "unauthorized"
    forbidden = "forbidden"
    not_found = "not_found"
    conflict_error = "conflict_error"
    totp_error = "totp_error"
    totp_invalid_backup_code_error = "totp_invalid_backup_code_error"
    totp_required_error = "totp_required_error"
    totp_exists_error = "totp_exists_error"
    registration_error = "registration_error"
    registrations_closed_error = "registrations_closed_error"
    password_reset_code_expired_error = "password_reset_code_expired_error"
    weak_password_error = "weak_password_error"
    invalid_password_error = "invalid_password_error"
    similar_passwords_error = "similar_passwords_error"
    document_is_too_large_error = "document_is_too_large_error"
    unsupported_document_format = "unsupported_document_format"
    confirmation_code_expired_error = "confirmation_code_expired_error"
    invite_expired_error = "invite_expired_error"
    internal_error = "internal_error"
    bad_request = "bad_request"


class BaseAPIException(Exception):
    default_code: ErrorCode = ErrorCode.internal_error
    default_message: str = "internal server error"
    default_status: int = status.HTTP_500_INTERNAL_SERVER_ERROR
    default_data: Any = None

    def __init__(
        self,
        code: Optional[ErrorCode] = None,
        status_code: Optional[int] = None,
        message: Optional[str] = None,
        data: Optional[Any] = None,
    ):
        self.code = code or self.default_code
        self.status = status_code or self.default_status
        self.message = message or self.default_message
        self.data = data or self.default_data

    def to_dict(self) -> Dict[Any, Any]:
        result = dict(
            code=self.code.value,
            status_code=self.status,
            message=self.message,
            data=self.data,
        )
        to_remove = set()
        for key, value in result.items():
            if value is None:
                to_remove.add(key)

        for key in to_remove:
            del result[key]

        return result


class InternalError(BaseAPIException):
    ...


class BadRequestError(BaseAPIException):
    default_code = ErrorCode.bad_request
    default_status = status.HTTP_400_BAD_REQUEST
    default_message = "bad request"


class TOTPInvalidBackupCodeError(BaseAPIException):
    default_code = ErrorCode.totp_invalid_backup_code_error
    default_status = status.HTTP_401_UNAUTHORIZED
    default_message = "invalid totp backup code"


class TOTPError(BaseAPIException):
    default_code = ErrorCode.totp_error
    default_status = status.HTTP_400_BAD_REQUEST
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


class ConflictError(BaseAPIException):
    default_code = ErrorCode.conflict_error
    default_status = status.HTTP_409_CONFLICT
    default_message = "conflict error"


class RegistrationError(BaseAPIException):
    default_code = ErrorCode.registration_error
    default_status = status.HTTP_409_CONFLICT
    default_message = "registration error"


class RegistrationsClosedError(BaseAPIException):
    default_code = ErrorCode.registrations_closed_error
    default_status = status.HTTP_400_BAD_REQUEST
    default_message = "registrations closed"


class ConfirmationCodeExpiredError(BaseAPIException):
    default_code = ErrorCode.confirmation_code_expired_error
    default_status = status.HTTP_409_CONFLICT
    default_message = "confirmation code has expired"


class InviteExpiredError(BaseAPIException):
    default_code = ErrorCode.invalid_password_error
    default_status = status.HTTP_409_CONFLICT
    default_message = "invite has expired"


class PasswordResetCodeExpiredError(BaseAPIException):
    default_code = ErrorCode.password_reset_code_expired_error
    default_status = status.HTTP_401_UNAUTHORIZED
    default_message = "password reset code hash expired"


class InvalidPasswordError(BaseAPIException):
    default_code = ErrorCode.invalid_password_error
    default_status = status.HTTP_403_FORBIDDEN
    default_message = "password is invalid"


class WeakPasswordError(BaseAPIException):
    default_code = ErrorCode.weak_password_error
    default_status = status.HTTP_400_BAD_REQUEST
    default_message = "password is weak"


class SimilarPasswordsError(BaseAPIException):
    default_code = ErrorCode.similar_passwords_error
    default_status = status.HTTP_400_BAD_REQUEST
    default_message = "passwords are too similar"


class UnsupportedDocumentError(BaseAPIException):
    default_code = ErrorCode.unsupported_document_format
    default_status = status.HTTP_415_UNSUPPORTED_MEDIA_TYPE
    default_message = "file upload is not permitted"


class TooLargeDocumentError(BaseAPIException):
    default_code = ErrorCode.document_is_too_large_error
    default_status = status.HTTP_413_REQUEST_ENTITY_TOO_LARGE
    default_message = "document is too large"
