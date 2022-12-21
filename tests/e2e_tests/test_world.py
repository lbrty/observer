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
