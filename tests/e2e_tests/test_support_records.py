import uuid

from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.types import BeneficiaryAge, SupportRecordSubject, SupportType
from observer.entities.permissions import NewPermission
from observer.schemas.support_records import (
    NewSupportRecordRequest,
    UpdateSupportRecordRequest,
)
from tests.helpers.auth import get_auth_tokens
from tests.helpers.crud import (
    create_permission,
    create_person,
    create_project,
    create_support_record,
)


async def test_create_support_record_works(
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
            can_create_projects=True,
            can_read_documents=True,
            can_read_personal_info=True,
            can_invite_members=True,
            project_id=default_project.id,
            user_id=consultant_user.id,
        ),
    )

    person = await create_person(app_context, default_project.id)
    payload = NewSupportRecordRequest(
        description="Buy clothes",
        type=SupportType.humanitarian,
        consultant_id=consultant_user.id,
        beneficiary_age=BeneficiaryAge.young_teen,
        record_for=SupportRecordSubject.person,
        owner_id=person.id,
        project_id=default_project.id,
    )
    resp = await authorized_client.post("/support-records", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED

    created_record = resp.json()
    record_id = created_record["id"]
    resp = await authorized_client.get(f"/support-records/{record_id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == created_record


async def test_create_support_relation_checks_work(
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
            can_create_projects=True,
            can_read_documents=True,
            can_read_personal_info=True,
            can_invite_members=True,
            project_id=default_project.id,
            user_id=consultant_user.id,
        ),
    )

    payload = NewSupportRecordRequest(
        description="Buy clothes",
        type=SupportType.humanitarian,
        consultant_id=consultant_user.id,
        beneficiary_age=BeneficiaryAge.young_teen,
        record_for=SupportRecordSubject.person,
        owner_id=uuid.uuid4(),
        project_id=default_project.id,
    )
    resp = await authorized_client.post("/support-records", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_404_NOT_FOUND

    payload = NewSupportRecordRequest(
        description="Buy clothes",
        type=SupportType.humanitarian,
        consultant_id=consultant_user.id,
        beneficiary_age=BeneficiaryAge.young_teen,
        record_for=SupportRecordSubject.pet,
        owner_id=uuid.uuid4(),
        project_id=default_project.id,
    )
    resp = await authorized_client.post("/support-records", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_404_NOT_FOUND


async def test_get_support_record_works(
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
            can_create_projects=True,
            can_read_documents=True,
            can_read_personal_info=True,
            can_invite_members=True,
            project_id=default_project.id,
            user_id=consultant_user.id,
        ),
    )

    person = await create_person(app_context, default_project.id)
    payload = NewSupportRecordRequest(
        description="Recover documents",
        type=SupportType.legal,
        consultant_id=consultant_user.id,
        beneficiary_age=BeneficiaryAge.young_adult,
        record_for=SupportRecordSubject.person,
        owner_id=person.id,
        project_id=default_project.id,
    )
    record = await create_support_record(app_context, payload)
    resp = await authorized_client.get(f"/support-records/{record.id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == jsonable_encoder(record)


async def test_update_support_record_works(
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
            can_create_projects=True,
            can_read_documents=True,
            can_read_personal_info=True,
            can_invite_members=True,
            project_id=default_project.id,
            user_id=consultant_user.id,
        ),
    )

    person = await create_person(app_context, default_project.id)
    payload = NewSupportRecordRequest(
        description="Buy clothes",
        type=SupportType.humanitarian,
        consultant_id=consultant_user.id,
        beneficiary_age=BeneficiaryAge.young_teen,
        record_for=SupportRecordSubject.person,
        owner_id=person.id,
        project_id=default_project.id,
    )
    record = await create_support_record(app_context, payload)
    updates = UpdateSupportRecordRequest(
        description="Buy clothes",
        type=SupportType.humanitarian,
        consultant_id=consultant_user.id,
        beneficiary_age=BeneficiaryAge.young_teen,
        record_for=SupportRecordSubject.person,
        owner_id=person.id,
        project_id=default_project.id,
    )
    resp = await authorized_client.put(f"/support-records/{record.id}", json=jsonable_encoder(updates))
    assert resp.status_code == status.HTTP_200_OK

    updated = resp.json()
    resp = await authorized_client.get(f"/support-records/{record.id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == updated


async def test_delete_support_record_works(
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
            can_create_projects=True,
            can_read_documents=True,
            can_read_personal_info=True,
            can_invite_members=True,
            project_id=default_project.id,
            user_id=consultant_user.id,
        ),
    )

    person = await create_person(app_context, default_project.id)
    payload = NewSupportRecordRequest(
        description="Buy clothes",
        type=SupportType.humanitarian,
        consultant_id=consultant_user.id,
        beneficiary_age=BeneficiaryAge.young_teen,
        record_for=SupportRecordSubject.person,
        owner_id=person.id,
        project_id=default_project.id,
    )
    record = await create_support_record(app_context, payload)
    resp = await authorized_client.delete(f"/support-records/{record.id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    resp = await authorized_client.get(f"/support-records/{record.id}")
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "message": "Support record not found",
        "status_code": 404,
    }


async def test_permission_checks_work(
    client,
    app_context,
    consultant_user,
    guest_user,
):
    project = await create_project(app_context, "Project Errors", description="Error catalog")
    person = await create_person(app_context, project.id)
    payload = NewSupportRecordRequest(
        description="Buy clothes",
        type=SupportType.humanitarian,
        consultant_id=consultant_user.id,
        beneficiary_age=BeneficiaryAge.young_teen,
        record_for=SupportRecordSubject.person,
        owner_id=person.id,
        project_id=project.id,
    )
    record = await create_support_record(app_context, payload)
    resp = await client.get(f"/support-records/{record.id}")
    assert resp.status_code == status.HTTP_401_UNAUTHORIZED

    auth_token = await get_auth_tokens(app_context, guest_user)
    resp = await client.get(
        f"/support-records/{record.id}",
        cookies=auth_token.dict(),
    )
    assert resp.status_code == status.HTTP_403_FORBIDDEN
    assert resp.json() == {"code": "unauthorized", "status_code": 403, "message": "Access forbidden"}
