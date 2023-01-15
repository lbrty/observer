from datetime import timedelta

AccessTokenKey: str = "access_token"
RefreshTokenKey: str = "refresh_token"

AccessTokenExpirationMinutes: int = 15
RefreshTokenExpirationMinutes: int = 180
AccessTokenExpirationDelta = timedelta(minutes=AccessTokenExpirationMinutes)
RefreshTokenExpirationDelta = timedelta(days=RefreshTokenExpirationMinutes)
