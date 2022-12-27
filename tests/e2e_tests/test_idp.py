from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.entities.permissions import NewPermission
from observer.schemas.idp import NewIDPRequest, UpdateIDPRequest
from tests.helpers.crud import create_permission, create_project


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
