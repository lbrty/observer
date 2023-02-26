from datetime import timedelta

from observer.settings import settings

AccessTokenKey: str = "access_token"
RefreshTokenKey: str = "refresh_token"

AccessTokenExpirationDelta = timedelta(minutes=settings.access_token_expiration_minutes)
RefreshTokenExpirationDelta = timedelta(days=settings.refresh_token_expiration_days)
