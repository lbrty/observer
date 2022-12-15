from starlette import status


async def test_create_project_works_as_expected(authorized_client, ensure_db, app_context, consultant_user):
    resp = await authorized_client.post(
        "/projects/",
        json=dict(
            name="Test Project",
            description="Project description",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    assert resp_json["name"] == "Test Project"
    assert resp_json["description"] == "Project description"
