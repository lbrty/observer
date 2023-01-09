import os
from datetime import date

from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.entities.permissions import NewPermission
from observer.schemas.idp import NewIDPRequest, UpdateIDPRequest
from observer.schemas.migration_history import (
    FullMigrationHistoryResponse,
    NewMigrationHistoryRequest,
)
from tests.helpers.crud import (
    create_city,
    create_country,
    create_permission,
    create_project,
    create_state,
)


async def test_create_idp_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
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
        tags=["one", "two"],
    )
    resp = await authorized_client.post("/idp/people", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED


async def test_get_idp_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
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
        phone_number_additional="+18181818",
        tags=["one", "two"],
    )
    resp = await authorized_client.post("/idp/people", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    resp_json["email"] = "********"
    resp_json["phone_number"] = "********"
    resp_json["phone_number_additional"] = "********"
    idp_id = resp_json["id"]
    resp = await authorized_client.get(f"/idp/people/{idp_id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == resp_json


async def test_get_idp_personal_info_works_as_expected(
    authorized_client,
    app_context,
    consultant_user,
):
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

    resp_json = resp.json()
    idp_id = resp_json["id"]
    resp = await authorized_client.get(f"/idp/people/{idp_id}/personal-info")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == {
        "email": "Full_Name@examples.com",
        "full_name": "Full Name",
        "phone_number": "+11111111",
        "phone_number_additional": "+18181818",
    }


async def test_update_idp_works_as_expected(
    authorized_client,
    app_context,
    consultant_user,
):
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

    resp_json = resp.json()
    idp_id = resp_json["id"]
    payload = UpdateIDPRequest(
        project_id=project.id,
        email="********",
        full_name="Full Name Updated",
        phone_number="+111111118888888",
        phone_number_additional="+48186818",
        tags=["one", "two", "three"],
    )
    resp = await authorized_client.put(f"/idp/people/{idp_id}", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_200_OK

    resp = await authorized_client.get(f"/idp/people/{idp_id}/personal-info")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == {
        "email": "Full_Name@examples.com",
        "full_name": "Full Name Updated",
        "phone_number": "+111111118888888",
        "phone_number_additional": "+48186818",
    }

    payload = UpdateIDPRequest(
        project_id=project.id,
        email="updated@email.com",
        full_name="Full Name Updated",
        phone_number="+111111118888888",
        phone_number_additional="+48186818",
        tags=["one", "two", "three"],
    )
    resp = await authorized_client.put(f"/idp/people/{idp_id}", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_200_OK

    resp = await authorized_client.get(f"/idp/people/{idp_id}/personal-info")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == {
        "email": "updated@email.com",
        "full_name": "Full Name Updated",
        "phone_number": "+111111118888888",
        "phone_number_additional": "+48186818",
    }

    resp = await authorized_client.get(f"/idp/people/{idp_id}")
    assert resp.status_code == status.HTTP_200_OK

    resp_json = resp.json()
    del resp_json["created_at"]
    del resp_json["updated_at"]
    assert resp_json == {
        "status": "registered",
        "reference_id": None,
        "email": "********",
        "full_name": "Full Name Updated",
        "birth_date": None,
        "notes": None,
        "phone_number": "********",
        "phone_number_additional": "********",
        "project_id": str(project.id),
        "category_id": None,
        "tags": ["one", "two", "three"],
        "id": idp_id,
        "external_id": None,
        "consultant_id": None,
    }


async def test_delete_idp_works_as_expected(
    authorized_client,
    app_context,
    consultant_user,
):
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

    resp_json = resp.json()
    idp_id = resp_json["id"]
    resp = await authorized_client.delete(f"/idp/people/{idp_id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    resp = await authorized_client.get(f"/idp/people/{idp_id}")
    assert resp.status_code == status.HTTP_404_NOT_FOUND


async def test_get_idp_migration_history_works_as_expected(authorized_client, app_context, consultant_user):
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
    country = await create_country(app_context, "Country 1", "c1")
    state = await create_state(app_context, "State 1", "s1", country.id)
    city_1 = await create_city(app_context, "City 1", "cty1", country.id, state.id)
    city_2 = await create_city(app_context, "City 2", "cty2", country.id, state.id)
    payload = NewMigrationHistoryRequest(
        idp_id=person_id,
        project_id=project.id,
        migration_date=date(year=2018, month=8, day=4),
        from_place_id=city_1.id,
        current_place_id=city_2.id,
    )
    resp = await authorized_client.post("/migrations", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED
    expected_response = FullMigrationHistoryResponse(**resp.json())
    expected_response.from_place = city_1
    expected_response.current_place = city_2

    resp = await authorized_client.get(f"/idp/people/{person_id}/migration-records")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == [jsonable_encoder(expected_response)]


async def test_upload_document_for_idp_works(
    authorized_client,
    app_context,
    default_idp,
    markdown_file,
    env_settings,
    fs_storage,
):
    files = {"file": ("readme.md", markdown_file, "text/markdown")}
    resp = await authorized_client.post(f"/idp/people/{default_idp.id}/document", files=files)
    assert resp.status_code == status.HTTP_201_CREATED
    document = await app_context.documents_service.get_document(resp.json()["id"])
    assert os.path.exists(document.path)


async def test_get_persons_documents_works(
    authorized_client,
    app_context,
    default_idp,
    markdown_file,
    textfile,
    env_settings,
    fs_storage,
):
    documents = []
    files = {"file": ("readme.md", markdown_file, "text/markdown")}
    resp = await authorized_client.post(f"/idp/people/{default_idp.id}/document", files=files)
    assert resp.status_code == status.HTTP_201_CREATED
    document = await app_context.documents_service.get_document(resp.json()["id"])
    assert os.path.exists(document.path)

    documents.append(resp.json())
    document = await app_context.documents_service.get_document(documents[0]["id"])
    assert os.path.exists(document.path)

    resp = await authorized_client.post(
        f"/idp/people/{default_idp.id}/document",
        files={"file": ("notes.txt", textfile, "text/plain")},
    )
    assert resp.status_code == status.HTTP_201_CREATED

    documents.append(resp.json())
    document = await app_context.documents_service.get_document(documents[1]["id"])
    assert os.path.exists(document.path)

    resp = await authorized_client.get(f"/idp/people/{default_idp.id}/documents")
    assert resp.json() == documents
