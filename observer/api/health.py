from fastapi import APIRouter

from observer.schemas.health import HealthResponse

router = APIRouter(prefix="/health")


@router.get(
    "",
    response_model=HealthResponse,
    tags=["health"],
)
async def health() -> HealthResponse:
    """Simple health endpoint"""
    return HealthResponse(status="ok")
