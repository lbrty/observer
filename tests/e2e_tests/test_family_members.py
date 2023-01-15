from datetime import date

from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.types import AgeGroup
from observer.schemas.family_members import NewFamilyMemberRequest


async def test_add_family_member_works(
    authorized_client,
    app_context,
    consultant_user,
    default_person,
    default_project,
):
    payload = NewFamilyMemberRequest(
        idp_id=default_person.id,
        age_group=AgeGroup.young_teen,
        project_id=default_project.id,
        migration_date=date(year=2018, month=8, day=4),
    )
    resp = await authorized_client.post(f"/people/{default_person.id}/family-members", json=jsonable_encoder(payload))
    assert resp.status_code == status.HTTP_201_CREATED


async def test_get_family_members_works(
    authorized_client,
    app_context,
    consultant_user,
    default_person,
    default_family,
):
    resp = await authorized_client.get(f"/people/{default_person.id}/family-members")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == [jsonable_encoder(default_family)]


async def test_update_family_member_works(
    authorized_client,
    app_context,
    consultant_user,
    default_person,
    default_family,
):
    updates = default_family.dict()
    updates["notes"] = "Fam"
    updates["age_group"] = AgeGroup.young_adult
    resp = await authorized_client.put(
        f"/people/{default_person.id}/family-members/{default_family.id}",
        json=jsonable_encoder(updates),
    )
    assert resp.status_code == status.HTTP_200_OK

    updated_member = resp.json()
    resp = await authorized_client.get(f"/people/{default_person.id}/family-members")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == [updated_member]


async def test_delete_family_member_works(
    authorized_client,
    app_context,
    consultant_user,
    default_person,
    default_family,
):
    resp = await authorized_client.get(f"/people/{default_person.id}/family-members")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == [jsonable_encoder(default_family)]

    resp = await authorized_client.delete(f"/people/{default_person.id}/family-members/{default_family.id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    resp = await authorized_client.get(f"/people/{default_person.id}/family-members")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == []
