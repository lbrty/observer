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
