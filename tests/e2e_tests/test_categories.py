import uuid

from starlette import status


async def test_create_vulnerability_category_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/idp/categories",
        json=dict(name="Vuln category"),
    )
    assert resp.status_code == status.HTTP_201_CREATED


async def test_get_vulnerability_category_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/idp/categories",
        json=dict(name="Vuln category"),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    category_id = resp_json["id"]
    resp = await authorized_client.get(f"/idp/categories/{category_id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == resp_json


async def test_get_unknown_vulnerability_category_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.get(f"/idp/categories/{uuid.uuid4()}")
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "message": "Category not found",
        "status_code": 404,
    }


async def test_get_vulnerability_categories_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    categories = []
    for n in range(10):
        resp = await authorized_client.post(
            "/idp/categories",
            json=dict(name=f"Vuln category #{n + 1}"),
        )
        assert resp.status_code == status.HTTP_201_CREATED
        categories.append(resp.json())

    resp = await authorized_client.get("/idp/categories")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == categories


async def test_filter_vulnerability_categories_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    for n in range(10):
        if n % 2 == 0:
            name = f"Catsegory #{n + 1}"
        else:
            name = f"Cowsegory #{n + 1}"

        resp = await authorized_client.post(
            "/idp/categories",
            json=dict(name=name),
        )
        assert resp.status_code == status.HTTP_201_CREATED

    resp = await authorized_client.get("/idp/categories?name=egory")
    assert resp.status_code == status.HTTP_200_OK
    assert len(resp.json()) == 10

    resp = await authorized_client.get("/idp/categories?name=cow")
    assert resp.status_code == status.HTTP_200_OK
    assert len(resp.json()) == 5

    resp = await authorized_client.get("/idp/categories?name=Cats")
    assert resp.status_code == status.HTTP_200_OK
    assert len(resp.json()) == 5


async def test_update_vulnerability_category_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/idp/categories",
        json=dict(name=f"Vulnerability category"),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    category_id = resp_json["id"]
    resp = await authorized_client.get(f"/idp/categories/{category_id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == resp_json


async def test_delete_vulnerability_category_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/idp/categories",
        json=dict(name=f"Vulnerability category"),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    category_id = resp_json["id"]
    resp = await authorized_client.delete(f"/idp/categories/{category_id}")
    assert resp.status_code == status.HTTP_204_NO_CONTENT
