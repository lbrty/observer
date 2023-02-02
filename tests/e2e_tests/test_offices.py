from fastapi.encoders import jsonable_encoder
from pytest import mark
from starlette import status

from observer.schemas.offices import NewOfficeRequest, UpdateOfficeRequest
from tests.helpers.auth import get_auth_tokens

OfficeNames = ["Bishkek", "Karakol", "Tokmok", "Kyiv", "Chornobaivka", "Tokmak"]


@mark.parametrize("office_name", OfficeNames)
async def test_create_office_works(office_name, app_context, client, admin_user):
    tokens = await get_auth_tokens(app_context, admin_user)
    payload = NewOfficeRequest(name=office_name)
    cookies = tokens.dict()
    resp = await client.post("/offices", json=jsonable_encoder(payload), cookies=cookies)
    assert resp.status_code == status.HTTP_201_CREATED

    office = resp.json()
    office_id = office["id"]
    resp = await client.get(f"/offices/{office_id}", cookies=cookies)
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == office


async def test_get_offices_works(app_context, client, admin_user):
    tokens = await get_auth_tokens(app_context, admin_user)
    cookies = tokens.dict()
    offices = []
    for office_name in OfficeNames:
        payload = NewOfficeRequest(name=office_name)
        resp = await client.post("/offices", json=jsonable_encoder(payload), cookies=cookies)
        assert resp.status_code == status.HTTP_201_CREATED
        offices.append(resp.json())

    resp = await client.get("/offices", cookies=cookies)
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == {
        "total": len(offices),
        "items": offices,
    }

    # Pagination
    resp = await client.get("/offices?offset=0&limit=2", cookies=cookies)
    assert resp.json() == {
        "total": len(offices),
        "items": offices[0:2],
    }

    # Filtering
    resp = await client.get("/offices?name=tok", cookies=cookies)
    assert resp.json() == {
        "total": len(offices),
        "items": [offices[2], offices[-1]],
    }

    # Lookup for Tokmok and Tokmak
    resp = await client.get("/offices?name=tok&offset=0&limit=2", cookies=cookies)
    assert resp.json() == {
        "total": len(offices),
        "items": [offices[2], offices[-1]],
    }


async def test_update_offices_works(app_context, client, admin_user):
    tokens = await get_auth_tokens(app_context, admin_user)
    cookies = tokens.dict()
    payload = NewOfficeRequest(name="Potsdam")
    resp = await client.post("/offices", json=jsonable_encoder(payload), cookies=cookies)
    assert resp.status_code == status.HTTP_201_CREATED

    office = resp.json()
    office_id = office["id"]
    payload = UpdateOfficeRequest(name="Potsdamsky Żabka")
    resp = await client.put(f"/offices/{office_id}", json=jsonable_encoder(payload), cookies=cookies)
    assert resp.status_code == status.HTTP_200_OK


async def test_delete_offices_works(app_context, client, admin_user):
    tokens = await get_auth_tokens(app_context, admin_user)
    cookies = tokens.dict()
    payload = NewOfficeRequest(name="Beijing")
    resp = await client.post("/offices", json=jsonable_encoder(payload), cookies=cookies)
    assert resp.status_code == status.HTTP_201_CREATED

    office = resp.json()
    office_id = office["id"]
    resp = await client.delete(f"/offices/{office_id}", cookies=cookies)
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    resp = await client.get(f"/offices/{office_id}", cookies=cookies)
    assert resp.status_code == status.HTTP_404_NOT_FOUND
