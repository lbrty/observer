from datetime import date

from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.schemas.migration_history import NewMigrationHistoryRequest


async def test_add_migration_history_works_as_expected(
    authorized_client,
    app_context,
    consultant_user,
    default_person,
    default_project,
):
    payload = NewMigrationHistoryRequest(
        idp_id=default_person.id,
        project_id=default_project.id,
        migration_date=date(year=2018, month=8, day=4),
    )
    resp = await authorized_client.post("/migrations", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED

    migration_record_id = resp.json()["id"]
    migration_record = await app_context.migrations_service.get_record(migration_record_id)
    assert jsonable_encoder(migration_record) == resp.json()


async def test_get_migration_history_works_as_expected(
    authorized_client,
    app_context,
    consultant_user,
    default_person,
    default_project,
):
    payload = NewMigrationHistoryRequest(
        idp_id=default_person.id,
        project_id=default_project.id,
        migration_date=date(year=2018, month=8, day=4),
    )
    resp = await authorized_client.post("/migrations", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED

    created_record = resp.json()
    migration_record_id = resp.json()["id"]
    resp = await authorized_client.get(f"/migrations/{migration_record_id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == created_record


async def test_delete_migration_history_works_as_expected(
    authorized_client,
    app_context,
    consultant_user,
    default_person,
    default_project,
):
    payload = NewMigrationHistoryRequest(
        idp_id=default_person.id,
        project_id=default_project.id,
        migration_date=date(year=2018, month=8, day=4),
    )
    resp = await authorized_client.post("/migrations", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED

    migration_record_id = resp.json()["id"]
    resp = await authorized_client.delete(f"/migrations/{migration_record_id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    resp = await authorized_client.get(f"/migrations/{migration_record_id}")
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "message": "Migration record not found",
        "status_code": 404,
    }
