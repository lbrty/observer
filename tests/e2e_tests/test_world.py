from starlette import status


async def test_create_country_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/world/countries",
        json=dict(
            name="Qyrgyzstan",
            code="qy",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    assert resp_json["name"] == "Qyrgyzstan"
    assert resp_json["code"] == "qy"


async def test_get_countries_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/world/countries",
        json=dict(
            name="Qyrgyzstan",
            code="qy",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp = await authorized_client.get("/world/countries")
    assert resp.status_code == status.HTTP_200_OK
    country = resp.json()[0]
    assert country["name"] == "Qyrgyzstan"
    assert country["code"] == "qy"


async def test_get_country_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/world/countries",
        json=dict(
            name="Qyrgyzstan",
            code="qy",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED
    country_id = resp.json()["id"]
    resp = await authorized_client.get(f"/world/countries/{country_id}")
    assert resp.status_code == status.HTTP_200_OK
    country = resp.json()
    assert country["name"] == "Qyrgyzstan"
    assert country["code"] == "qy"


async def test_update_country_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/world/countries",
        json=dict(
            name="Qyrgyzstan",
            code="qy",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED
    country_id = resp.json()["id"]
    resp = await authorized_client.put(
        f"/world/countries/{country_id}",
        json=dict(
            name="Qyrgyz Ordosu",
            code="qy",
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    resp = await authorized_client.get(f"/world/countries/{country_id}")
    assert resp.status_code == status.HTTP_200_OK
    country = resp.json()
    assert country["name"] == "Qyrgyz Ordosu"
    assert country["code"] == "qy"


async def test_delete_country_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    countries = []
    for n in range(10):
        resp = await authorized_client.post(
            "/world/countries",
            json=dict(
                name=f"Qyrgyzstan {n}",
                code=f"qy{n}",
            ),
        )
        assert resp.status_code == status.HTTP_201_CREATED
        countries.append(resp.json())
    assert len(countries) == 10

    await authorized_client.delete(f"/world/countries/{countries[0]['id']}")
    await authorized_client.delete(f"/world/countries/{countries[1]['id']}")

    resp = await authorized_client.get("/world/countries")
    assert resp.status_code == status.HTTP_200_OK
    assert len(resp.json()) == 8
