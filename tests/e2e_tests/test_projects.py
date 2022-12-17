from starlette import status

from tests.helpers.auth import get_auth_tokens


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


async def test_update_project_works_as_expected_for_members(authorized_client, ensure_db, app_context, consultant_user):
    resp = await authorized_client.post(
        "/projects/",
        json=dict(
            name="Test Project",
            description="Project description",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    resp = await authorized_client.put(
        f"/projects/{resp_json['id']}",
        json=dict(
            name="Test Project Updated",
            description="Project description updated",
        ),
    )
    assert resp.status_code == status.HTTP_200_OK
    project_id = resp_json["id"]
    assert resp.json() == dict(
        id=resp_json["id"],
        name="Test Project Updated",
        description="Project description updated",
    )
    audit_log = await app_context.audit_service.find_by_ref(
        "source=service:projects,endpoint=update_project,action=update:project,"
        f"project_id={project_id},ref_id={consultant_user.ref_id}",
    )

    assert audit_log.data == {
        "new_project": {
            "name": "Test Project Updated",
            "description": "Project description updated",
        },
        "old_project": {
            "name": "Test Project",
            "description": "Project description",
        },
    }


async def test_update_project_works_as_expected_for_admins(
    authorized_client,
    client,
    ensure_db,
    app_context,
    admin_user,
):
    auth_token = await get_auth_tokens(app_context, admin_user)
    resp = await authorized_client.post(
        "/projects/",
        json=dict(
            name="Test Project",
            description="Project description",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    resp = await client.put(
        f"/projects/{resp_json['id']}",
        json=dict(
            name="Test Project Updated",
            description="Project description updated",
        ),
        cookies=auth_token.dict(),
    )

    assert resp.status_code == status.HTTP_200_OK
    project_id = resp_json["id"]
    assert resp.json() == dict(
        id=project_id,
        name="Test Project Updated",
        description="Project description updated",
    )
    audit_log = await app_context.audit_service.find_by_ref(
        "source=service:projects,endpoint=update_project,action=update:project,"
        f"project_id={project_id},ref_id={admin_user.ref_id}",
    )

    assert audit_log.data == {
        "new_project": {
            "name": "Test Project Updated",
            "description": "Project description updated",
        },
        "old_project": {
            "name": "Test Project",
            "description": "Project description",
        },
    }


async def test_update_project_works_as_expected_for_users_without_permissions(
    authorized_client,
    client,
    ensure_db,
    app_context,
    guest_user,
):
    auth_token = await get_auth_tokens(app_context, guest_user)
    resp = await authorized_client.post(
        "/projects/",
        json=dict(
            name="Test Project",
            description="Project description",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    resp = await client.put(
        f"/projects/{resp_json['id']}",
        json=dict(
            name="Test Project Updated",
            description="Project description updated",
        ),
        cookies=auth_token.dict(),
    )

    assert resp.status_code == status.HTTP_403_FORBIDDEN
    assert resp.json() == {
        "code": "unauthorized",
        "status_code": 403,
        "message": "Yoo can not view this project",
        "data": None,
    }


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
    project_id = resp_json["id"]
    resp = await authorized_client.get(f"/projects/{project_id}/members")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == {
        "items": [
            {
                "ref_id": "ref-consultant-1",
                "is_active": True,
                "full_name": "full name",
                "role": "consultant",
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
