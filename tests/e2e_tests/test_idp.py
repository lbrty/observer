from datetime import datetime

from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.entities.permissions import NewPermission
from observer.schemas.idp import NewIDPRequest
from tests.helpers.crud import create_permission, create_project


async def test_create_idp_works_as_expected(authorized_client, ensure_db, app_context, consultant_user):
    project = await create_project(app_context, "test project", "test description")
    await create_permission(
        app_context,
        NewPermission(
            can_create=True,
            can_read=True,
            can_update=True,
            can_delete=True,
            can_create_projects=True,
            can_read_documents=True,
            can_read_personal_info=True,
            can_invite_members=True,
            project_id=project.id,
            user_id=consultant_user.id,
        ),
    )

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


async def test_get_idp_works_as_expected(authorized_client, ensure_db, app_context, consultant_user):
    project = await create_project(app_context, "test project", "test description")
    await create_permission(
        app_context,
        NewPermission(
            can_create=True,
            can_read=True,
            can_update=True,
            can_delete=False,
            can_create_projects=True,
            can_read_documents=False,
            can_read_personal_info=False,
            can_invite_members=False,
            project_id=project.id,
            user_id=consultant_user.id,
        ),
    )
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

    resp_json = resp.json()
    idp_id = resp_json["id"]
    resp = await authorized_client.get(f"/idp/people/{idp_id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == resp_json
