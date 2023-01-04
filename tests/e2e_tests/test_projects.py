import uuid

from starlette import status

from tests.helpers.auth import get_auth_tokens


async def test_create_project_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
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


async def test_get_project_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
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
        owner_id=str(consultant_user.id),
    )


async def test_update_project_works_as_expected_for_members(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
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
    resp = await authorized_client.put(
        f"/projects/{project_id}",
        json=dict(
            name="Test Project Updated",
            description="Project description updated",
        ),
    )
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == dict(
        id=resp_json["id"],
        name="Test Project Updated",
        description="Project description updated",
        owner_id=str(consultant_user.id),
    )
    audit_log = await app_context.audit_service.find_by_ref(
        "source=service:projects,endpoint=update_project,action=update:project,"
        f"project_id={project_id},ref_id={consultant_user.ref_id}",
    )

    assert audit_log.data == {
        "new_project": {
            "id": project_id,
            "name": "Test Project Updated",
            "description": "Project description updated",
            "owner_id": str(consultant_user.id),
        },
        "old_project": {
            "id": project_id,
            "name": "Test Project",
            "description": "Project description",
            "owner_id": str(consultant_user.id),
        },
    }


async def test_update_project_works_as_expected_for_admins(
    authorized_client,
    client,
    ensure_db,
    app_context,
    admin_user,
    consultant_user,
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
        owner_id=str(consultant_user.id),
    )
    audit_log = await app_context.audit_service.find_by_ref(
        "source=service:projects,endpoint=update_project,action=update:project,"
        f"project_id={project_id},ref_id={admin_user.ref_id}",
    )

    assert audit_log.data == {
        "new_project": {
            "id": project_id,
            "name": "Test Project Updated",
            "description": "Project description updated",
            "owner_id": str(consultant_user.id),
        },
        "old_project": {
            "id": project_id,
            "name": "Test Project",
            "description": "Project description",
            "owner_id": str(consultant_user.id),
        },
    }


async def test_delete_project_works_as_expected_for_members(authorized_client, ensure_db, app_context, consultant_user):
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
    resp = await authorized_client.delete(f"/projects/{project_id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    project_id = resp_json["id"]
    audit_log = await app_context.audit_service.find_by_ref(
        "source=service:projects,endpoint=delete_project,action=delete:project,"
        f"project_id={project_id},ref_id={consultant_user.ref_id}",
    )

    assert audit_log.data is None


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
        "message": "You can not update this project",
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
    permission = await app_context.repos.permissions.find(project_id, consultant_user.id)
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
                    "id": str(permission.id),
                    "can_create": True,
                    "can_read": True,
                    "can_update": True,
                    "can_delete": True,
                    "can_create_projects": True,
                    "can_read_documents": True,
                    "can_read_personal_info": True,
                    "can_invite_members": True,
                    "project_id": project_id,
                    "user_id": str(consultant_user.id),
                },
            }
        ],
    }


async def test_add_project_member_works_as_expected(
    authorized_client, ensure_db, app_context, consultant_user, guest_user
):
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
    resp = await authorized_client.post(
        f"/projects/{project_id}/members",
        json=dict(
            can_create=False,
            can_read=True,
            can_update=True,
            can_delete=False,
            can_create_projects=False,
            can_read_documents=False,
            can_read_personal_info=False,
            can_invite_members=False,
            project_id=project_id,
            user_id=str(guest_user.id),
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    consultant_permission = await app_context.repos.permissions.find(project_id, consultant_user.id)
    guest_permission = await app_context.repos.permissions.find(project_id, guest_user.id)
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
                    "id": str(consultant_permission.id),
                    "can_create": True,
                    "can_read": True,
                    "can_update": True,
                    "can_delete": True,
                    "can_create_projects": True,
                    "can_read_documents": True,
                    "can_read_personal_info": True,
                    "can_invite_members": True,
                    "project_id": project_id,
                    "user_id": str(consultant_user.id),
                },
            },
            {
                "ref_id": "ref-guest-1",
                "is_active": True,
                "full_name": "full name",
                "role": "guest",
                "permissions": {
                    "id": str(guest_permission.id),
                    "can_create": False,
                    "can_read": True,
                    "can_update": True,
                    "can_delete": False,
                    "can_create_projects": False,
                    "can_read_documents": False,
                    "can_read_personal_info": False,
                    "can_invite_members": False,
                    "project_id": project_id,
                    "user_id": str(guest_user.id),
                },
            },
        ]
    }


async def test_add_project_member_works_as_expected_for_users_without_permissions(
    authorized_client, ensure_db, app_context, consultant_user, guest_user
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
    project_id = resp_json["id"]
    resp = await authorized_client.post(
        f"/projects/{project_id}/members",
        json=dict(
            can_create=False,
            can_read=True,
            can_update=True,
            can_delete=False,
            can_create_projects=False,
            can_read_documents=False,
            can_read_personal_info=False,
            can_invite_members=False,
            project_id=project_id,
            user_id=str(guest_user.id),
        ),
        cookies=auth_token.dict(),
    )
    assert resp.status_code == status.HTTP_403_FORBIDDEN
    assert resp.json() == {
        "code": "unauthorized",
        "status_code": 403,
        "message": "You can not invite members in this project",
        "data": None,
    }


async def test_add_project_member_works_as_expected_if_given_user_does_not_exist(
    authorized_client, ensure_db, app_context, consultant_user, guest_user
):
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
    resp = await authorized_client.post(
        f"/projects/{project_id}/members",
        json=dict(
            can_create=False,
            can_read=True,
            can_update=True,
            can_delete=False,
            can_create_projects=False,
            can_read_documents=False,
            can_read_personal_info=False,
            can_invite_members=False,
            project_id=project_id,
            user_id=str(uuid.uuid4()),
        ),
    )
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "status_code": 404,
        "message": "User not found",
        "data": None,
    }


async def test_add_project_member_works_as_expected_if_project_ids_mismatch(
    authorized_client, ensure_db, app_context, consultant_user, guest_user
):
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
    resp = await authorized_client.post(
        f"/projects/{project_id}/members",
        json=dict(
            can_create=False,
            can_read=True,
            can_update=True,
            can_delete=False,
            can_create_projects=False,
            can_read_documents=False,
            can_read_personal_info=False,
            can_invite_members=False,
            project_id=str(uuid.uuid4()),
            user_id=str(guest_user.id),
        ),
    )
    assert resp.status_code == status.HTTP_409_CONFLICT
    assert resp.json() == {
        "code": "conflict_error",
        "status_code": 409,
        "message": "Project ids in path and payload differ",
        "data": None,
    }


async def test_delete_project_member_works_as_expected(
    authorized_client, ensure_db, app_context, consultant_user, guest_user
):
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
    resp = await authorized_client.post(
        f"/projects/{project_id}/members",
        json=dict(
            can_create=False,
            can_read=True,
            can_update=True,
            can_delete=False,
            can_create_projects=False,
            can_read_documents=False,
            can_read_personal_info=False,
            can_invite_members=False,
            project_id=project_id,
            user_id=str(guest_user.id),
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    consultant_permission = await app_context.repos.permissions.find(project_id, consultant_user.id)
    resp = await authorized_client.delete(f"/projects/{project_id}/members/{guest_user.id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT

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
                    "id": str(consultant_permission.id),
                    "can_create": True,
                    "can_read": True,
                    "can_update": True,
                    "can_delete": True,
                    "can_create_projects": True,
                    "can_read_documents": True,
                    "can_read_personal_info": True,
                    "can_invite_members": True,
                    "project_id": project_id,
                    "user_id": str(consultant_user.id),
                },
            }
        ]
    }


async def test_update_project_member_permissions_works_as_expected(
    authorized_client, ensure_db, app_context, consultant_user, guest_user
):
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
    consultant_permission = await app_context.repos.permissions.find(project_id, consultant_user.id)
    resp = await authorized_client.put(
        f"/projects/{project_id}/members/{consultant_user.id}",
        json={
            "can_create": True,
            "can_read": True,
            "can_update": True,
            "can_delete": True,
            "can_create_projects": False,
            "can_read_documents": False,
            "can_read_personal_info": False,
            "can_invite_members": False,
        },
    )
    assert resp.status_code == status.HTTP_200_OK

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
                    "id": str(consultant_permission.id),
                    "can_create": True,
                    "can_read": True,
                    "can_update": True,
                    "can_delete": True,
                    "can_create_projects": False,
                    "can_read_documents": False,
                    "can_read_personal_info": False,
                    "can_invite_members": False,
                    "project_id": project_id,
                    "user_id": str(consultant_user.id),
                },
            }
        ]
    }
