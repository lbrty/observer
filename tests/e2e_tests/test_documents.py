import os

from starlette import status


async def test_download_document_works_as_expected(
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

    doc_id = resp.json()["id"]
    resp = await authorized_client.get(f"/documents/stream/{doc_id}")
    assert resp.status_code == status.HTTP_200_OK
    markdown_file.seek(0)
    assert resp.content == markdown_file.read()
    assert resp.headers == {
        "content-disposition": 'text/markdown; filename="readme.md"',
        "content-type": "text/markdown; charset=utf-8",
    }


async def test_delete_document_works_as_expected(
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

    doc_id = resp.json()["id"]
    folder_path = os.path.join(env_settings.storage_root, env_settings.documents_path, str(new_pet.id))
    documents = await app_context.storage.ls(folder_path)
    assert len(documents) == 1

    resp = await authorized_client.delete(f"/documents/{doc_id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    documents = await app_context.storage.ls(folder_path)
    assert len(documents) == 0
