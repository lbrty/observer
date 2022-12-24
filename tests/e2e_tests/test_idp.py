from datetime import datetime

from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.schemas.idp import NewIDPRequest
from tests.helpers.crud import create_project


async def test_create_idp_works_as_expected(authorized_client, ensure_db, app_context, consultant_user):
    project = await create_project(app_context, "test project", "test description")
    payload = NewIDPRequest(
        project_id=project.id,
        email="Full_Name@examples.com",
        full_name="Full Name",
        phone_number="+11111111",
        migration_date=datetime.today(),
        tags=["one", "two"],
    )
    resp = await authorized_client.post("/idp/people", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED
