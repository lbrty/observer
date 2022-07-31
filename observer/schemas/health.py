from pydantic.main import BaseModel


class HealthResponse(BaseModel):
    status: str
