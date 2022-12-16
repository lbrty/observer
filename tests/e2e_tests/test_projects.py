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


async def test_get_project_works_as_expected(authorized_client, ensure_db, app_context, consultant_user):
    resp = await authorized_client.post(
        "/projects/",
        json=dict(
            name="Test Project",
            description="Project description",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    resp = await authorized_client.get(f"/projects/{resp_json['id']}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == dict(
        id=resp_json["id"],
        name="Test Project",
        description="Project description",
    )


async def test_get_project_members_works_as_expected(authorized_client, ensure_db, app_context, consultant_user):
    resp = await authorized_client.post(
        "/projects/",
        json=dict(
            name="Test Project",
            description="Project description",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED
    resp_json = resp.json()
    resp = await authorized_client.get(f"/projects/{resp_json['id']}/members")
    assert resp.status_code == status.HTTP_200_OK
    member = {
        "email": "consultant-1@example.com",
        "full_name": "full name",
        "role": "consultant",
        "id": str(consultant_user.id),
        "ref_id": "ref-consultant-1",
        "is_active": True,
        "is_confirmed": True,
        "mfa_enabled": False,
    }
    assert resp.json() == {
        "items": [
            {
                "user": member,
                "permissions": {
                    "can_create": True,
                    "can_read": True,
                    "can_update": True,
                    "can_delete": True,
                    "can_create_projects": True,
                    "can_read_documents": True,
                    "can_read_personal_info": True,
                },
            }
        ],
    }
