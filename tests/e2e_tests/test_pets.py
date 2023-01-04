import shutil
from io import BytesIO

from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.types import PetStatus
from observer.entities.permissions import NewPermission
from observer.schemas.idp import NewIDPRequest
from tests.helpers.crud import create_permission, create_pet, create_project


async def test_upload_document_for_pet_works(authorized_client, app_context, consultant_user):
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

    owner_id = resp.json()["id"]
    pet = await create_pet(
        app_context,
        "Jack",
        PetStatus.needs_shelter,
        "reg-id",
        project.id,
        owner_id,
    )
    # TODO: use tempdir
    fp = BytesIO()
    fp.write(b"BABA BLACK SHEEP")
    fp.seek(0)
    files = {"file": ("readme.md", fp, "text/markdown")}
    resp = await authorized_client.post(f"/pets/{pet.id}/document", files=files)
    assert resp.status_code == status.HTTP_201_CREATED
    contents = open(f"{app_context.storage.root}/readme.md").read()
    assert contents == "BABA BLACK SHEEP"
    shutil.rmtree(app_context.storage.root)
