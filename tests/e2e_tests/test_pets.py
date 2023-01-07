import os.path

from starlette import status


async def test_delete_pets_deletes_all_related_documents(
    authorized_client,
    app_context,
    consultant_user,
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


async def test_upload_document_for_pet_works(
    authorized_client,
    app_context,
    consultant_user,
    new_pet,
    markdown_file,
    env_settings,
    fs_storage,
):
    files = {"file": ("readme.md", markdown_file, "text/markdown")}
    resp = await authorized_client.post(f"/pets/{new_pet.id}/document", files=files)
    assert resp.status_code == status.HTTP_201_CREATED
    file_path = os.path.join(app_context.storage.root, env_settings.documents_path, str(new_pet.id), "readme.md")
    assert os.path.exists(file_path)
