import asyncio
import os.path
import uuid

import pytest
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.api.exceptions import NotFoundError
from observer.common.types import PetStatus
from observer.entities.permissions import NewPermission
from observer.schemas.pets import NewPetRequest, UpdatePetRequest
from tests.helpers.crud import create_permission, create_person


async def test_create_pet_works(
    authorized_client,
    app_context,
    default_project,
    consultant_user,
):
    await create_permission(
        app_context,
        NewPermission(
            can_create=True,
            can_read=True,
            can_update=True,
            can_delete=True,
            can_read_documents=True,
            can_read_personal_info=True,
            can_invite_members=True,
            project_id=default_project.id,
            user_id=consultant_user.id,
        ),
    )

    person = await create_person(app_context, default_project.id)
    payload = NewPetRequest(
        name="Moon",
        notes="Good doggo",
        status=PetStatus.owner_found,
        registration_id="chuy-1-111-11",
        owner_id=person.id,
        project_id=default_project.id,
    )
    resp = await authorized_client.post("/pets", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED

    pet = resp.json()
    resp = await authorized_client.get(f"/pets/{pet['id']}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == pet


async def test_upload_document_for_pet_works(
    authorized_client,
    app_context,
    new_pet,
    markdown_file,
    env_settings,
    fs_storage,
):
    files = {"file": ("readme.md", markdown_file, "text/markdown")}
    resp = await authorized_client.post(f"/pets/{new_pet.id}/document", files=files)
    assert resp.status_code == status.HTTP_201_CREATED
    document = await app_context.documents_service.get_document(resp.json()["id"])
    assert os.path.exists(document.path)


async def test_get_pets_documents_works(
    authorized_client,
    app_context,
    new_pet,
    markdown_file,
    textfile,
    env_settings,
    fs_storage,
):
    documents = []
    resp = await authorized_client.post(
        f"/pets/{new_pet.id}/document",
        files={"file": ("readme.md", markdown_file, "text/markdown")},
    )
    assert resp.status_code == status.HTTP_201_CREATED

    documents.append(resp.json())
    document = await app_context.documents_service.get_document(documents[0]["id"])
    assert os.path.exists(document.path)

    resp = await authorized_client.post(
        f"/pets/{new_pet.id}/document",
        files={"file": ("notes.txt", textfile, "text/plain")},
    )
    assert resp.status_code == status.HTTP_201_CREATED

    documents.append(resp.json())
    document = await app_context.documents_service.get_document(documents[1]["id"])
    assert os.path.exists(document.path)

    resp = await authorized_client.get(f"/pets/{new_pet.id}/documents")
    assert resp.json() == documents


async def test_get_pets_documents_with_remote_storage_works(
    authorized_client,
    app_context,
    new_pet,
    markdown_file,
    textfile,
    env_settings,
    aws_credentials,
    create_bucket,
    s3_storage,
    s3_client,
):
    documents = []
    resp = await authorized_client.post(
        f"/pets/{new_pet.id}/document",
        files={"file": ("readme.md", markdown_file, "text/markdown")},
    )
    assert resp.status_code == status.HTTP_201_CREATED
    documents.append(resp.json())

    resp = await authorized_client.post(
        f"/pets/{new_pet.id}/document",
        files={"file": ("notes.txt", textfile, "text/plain")},
    )
    assert resp.status_code == status.HTTP_201_CREATED
    documents.append(resp.json())

    resp = await authorized_client.get(f"/pets/{new_pet.id}/documents")
    assert resp.json() == documents

    folder_path = os.path.join(env_settings.storage_root, env_settings.documents_path, str(new_pet.id))
    documents = await app_context.storage.ls(folder_path)
    assert len(documents) == 2


async def test_update_pets_works(
    authorized_client,
    app_context,
    new_pet,
    env_settings,
):
    updates = UpdatePetRequest(**new_pet.dict())
    updates.name = "Toto"
    updates.status = PetStatus.owner_found
    resp = await authorized_client.put(f"/pets/{new_pet.id}", json=jsonable_encoder(updates, exclude={"id"}))
    assert resp.status_code == status.HTTP_200_OK
    resp_json = resp.json()
    resp_json["name"] = updates.name
    resp_json["status"] = PetStatus.owner_found.value

    resp = await authorized_client.put(f"/pets/{uuid.uuid4()}", json=jsonable_encoder(updates, exclude={"id"}))
    assert resp.status_code == status.HTTP_404_NOT_FOUND

    with pytest.raises(NotFoundError):
        await app_context.pets_service.update_pet(uuid.uuid4(), updates)


async def test_delete_pets_deletes_all_related_documents_from_remote_storage(
    authorized_client,
    app_context,
    env_settings,
    new_pet,
    markdown_file,
    aws_credentials,
    s3_storage,
    s3_client,
):
    files = {"file": ("readme.md", markdown_file, "text/markdown")}
    resp = await authorized_client.post(f"/pets/{new_pet.id}/document", files=files)
    assert resp.status_code == status.HTTP_201_CREATED

    resp = await authorized_client.delete(f"/pets/{new_pet.id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT
    folder_path = os.path.join(env_settings.documents_path, str(new_pet.id))
    documents = await app_context.storage.ls(folder_path)
    await asyncio.sleep(1)
    assert len(documents) == 0


async def test_delete_pets_deletes_all_related_documents(
    authorized_client,
    app_context,
    env_settings,
    new_pet,
    markdown_file,
    fs_storage,
):
    files = {"file": ("readme.md", markdown_file, "text/markdown")}
    resp = await authorized_client.post(f"/pets/{new_pet.id}/document", files=files)
    assert resp.status_code == status.HTTP_201_CREATED

    resp = await authorized_client.delete(f"/pets/{new_pet.id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT
    file_path = os.path.join(app_context.storage.root, env_settings.documents_path, str(new_pet.id))
    assert os.path.exists(file_path) is False
