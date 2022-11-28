from pydantic import EmailStr, Field, SecretStr

from observer.schemas.base import SchemaBase


class TokenResponse(SchemaBase):
    access_token: str = Field(..., description="JWT access token")
    refresh_token: str = Field(..., description="JWT refresh token")


class LoginPayload(SchemaBase):
    email: EmailStr = Field(..., description="E-mail address to login")
    password: SecretStr = Field(..., description="Password to login")


class RegistrationPayload(SchemaBase):
    email: EmailStr = Field(..., description="E-mail address of a user")
    password: SecretStr = Field(..., description="Password which user has provided")
