from datetime import date

from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.entities.permissions import NewPermission
from observer.schemas.idp import NewIDPRequest
from observer.schemas.migration_history import NewMigrationHistoryRequest
from tests.helpers.crud import create_permission, create_project


async def test_add_migration_history_works_as_expected(authorized_client, app_context, consultant_user):
    project = await create_project(app_context, "test project", "test description")
    await create_permission(
        app_context,
        NewPermission(
            can_create=True,
            can_read=True,
            can_update=True,
            can_delete=True,
            can_create_projects=True,
            can_read_documents=False,
            can_read_personal_info=True,
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
        phone_number_additional="+18181818",
        tags=["one", "two"],
    )
    resp = await authorized_client.post("/idp/people", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED
    person_id = resp.json()["id"]
    payload = NewMigrationHistoryRequest(
        idp_id=person_id,
        project_id=project.id,
        migration_date=date(year=2018, month=8, day=4),
    )
    resp = await authorized_client.post("/migrations", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED


async def test_get_migration_history_works_as_expected(authorized_client, app_context, consultant_user):
    pass


async def test_delete_migration_history_works_as_expected(authorized_client, app_context, consultant_user):
    pass
