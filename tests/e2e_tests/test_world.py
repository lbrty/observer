import uuid

from starlette import status

from observer.common.types import PlaceType


# Countries
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


async def test_update_unknown_country_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.put(
        f"/world/countries/{uuid.uuid4()}",
        json=dict(name="X", code="Y"),
    )
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "data": None,
        "message": "Country not found",
        "status_code": 404,
    }


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


async def test_delete_unknown_country_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.delete(f"/world/countries/{uuid.uuid4()}")
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "data": None,
        "message": "Country not found",
        "status_code": 404,
    }


# States
async def test_create_state_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
):
    resp = await authorized_client.post(
        "/world/states",
        json=dict(
            name="Qoçqor",
            code="qr",
            country_id=str(default_country.id),
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    assert resp_json["name"] == "Qoçqor"
    assert resp_json["code"] == "qr"
    assert resp_json["country_id"] == str(default_country.id)


async def test_get_states_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
):
    states = []
    for n in range(10):
        resp = await authorized_client.post(
            "/world/states",
            json=dict(
                name=f"Qoçqor {n}",
                code=f"qr{n}",
                country_id=str(default_country.id),
            ),
        )
        assert resp.status_code == status.HTTP_201_CREATED
        states.append(resp.json())

    assert len(states) == 10

    resp = await authorized_client.get("/world/states")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == states


async def test_get_states_with_filters_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
):
    states = []
    for n in range(10):
        resp = await authorized_client.post(
            "/world/states",
            json=dict(
                name=f"Qoçqor {n}",
                code=f"qr{n}",
                country_id=str(default_country.id),
            ),
        )
        assert resp.status_code == status.HTTP_201_CREATED
        states.append(resp.json())

    assert len(states) == 10

    resp = await authorized_client.get(
        "/world/states",
        params=dict(name="Qoçqor 1"),
    )
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json()[0] == states[1]


async def test_get_state_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
):
    resp = await authorized_client.post(
        "/world/states",
        json=dict(
            name="Qoçqor",
            code="qr",
            country_id=str(default_country.id),
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED
    state_id = resp.json()["id"]
    resp = await authorized_client.get(f"/world/states/{state_id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == dict(
        id=state_id,
        name="Qoçqor",
        code="qr",
        country_id=str(default_country.id),
    )


async def test_update_state_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
):
    resp = await authorized_client.post(
        "/world/states",
        json=dict(
            name="Qoçqor",
            code="qr",
            country_id=str(default_country.id),
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED
    state_id = resp.json()["id"]
    resp = await authorized_client.put(
        f"/world/states/{state_id}",
        json=dict(
            name="Balykçy",
            code="bc",
            country_id=str(default_country.id),
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    resp = await authorized_client.get(f"/world/states/{state_id}")
    assert resp.status_code == status.HTTP_200_OK
    country = resp.json()
    assert country["name"] == "Balykçy"
    assert country["code"] == "bc"


async def test_update_unknown_state_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.put(
        f"/world/states/{uuid.uuid4()}",
        json=dict(
            name="X",
            code="Y",
            country_id=str(uuid.uuid4()),
        ),
    )
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "data": None,
        "message": "Country not found",
        "status_code": 404,
    }


async def test_delete_state_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
):
    countries = []
    for n in range(10):
        resp = await authorized_client.post(
            "/world/states",
            json=dict(
                name=f"Qyrgyzstan {n}",
                code=f"qy{n}",
                country_id=str(default_country.id),
            ),
        )
        assert resp.status_code == status.HTTP_201_CREATED
        countries.append(resp.json())

    assert len(countries) == 10

    await authorized_client.delete(f"/world/states/{countries[0]['id']}")
    await authorized_client.delete(f"/world/states/{countries[1]['id']}")

    resp = await authorized_client.get("/world/states")
    assert resp.status_code == status.HTTP_200_OK
    assert len(resp.json()) == 8


# Places
async def test_create_place_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
    default_state,
):
    resp = await authorized_client.post(
        "/world/places",
        json=dict(
            name="Qoçqor",
            code="qr",
            place_type=PlaceType.town.value,
            country_id=str(default_country.id),
            state_id=str(default_state.id),
        ),
    )
    place_id = resp.json()["id"]
    assert resp.status_code == status.HTTP_201_CREATED
    assert resp.json() == dict(
        id=place_id,
        name="Qoçqor",
        code="qr",
        place_type=PlaceType.town.value,
        country_id=str(default_country.id),
        state_id=str(default_state.id),
    )


async def test_get_places_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
    default_state,
):
    places = []
    for n in range(10):
        resp = await authorized_client.post(
            "/world/places",
            json=dict(
                name=f"Bişkek {n}",
                code=f"bi{n}",
                place_type=PlaceType.city.value,
                country_id=str(default_country.id),
                state_id=str(default_state.id),
            ),
        )
        assert resp.status_code == status.HTTP_201_CREATED
        places.append(resp.json())

    resp = await authorized_client.get("/world/places")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == places


async def test_get_place_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
    default_state,
):
    resp = await authorized_client.post(
        "/world/places",
        json=dict(
            name="Calal-Abad",
            code="ca",
            place_type=PlaceType.city.value,
            country_id=str(default_country.id),
            state_id=str(default_state.id),
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    place_id = resp.json()["id"]
    resp = await authorized_client.get(f"/world/places/{place_id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == dict(
        id=place_id,
        name="Calal-Abad",
        code="ca",
        place_type=PlaceType.city.value,
        country_id=str(default_country.id),
        state_id=str(default_state.id),
    )


async def test_update_place_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
    default_state,
):
    resp = await authorized_client.post(
        "/world/places",
        json=dict(
            name="Oş",
            code="os",
            place_type=PlaceType.city.value,
            country_id=str(default_country.id),
            state_id=str(default_state.id),
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED
    place_id = resp.json()["id"]
    resp = await authorized_client.put(
        f"/world/places/{place_id}",
        json=dict(
            name="Oş",
            code="oh",
            place_type=PlaceType.city.value,
            country_id=str(default_country.id),
            state_id=str(default_state.id),
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    resp = await authorized_client.get(f"/world/places/{place_id}")
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == dict(
        id=place_id,
        name="Oş",
        code="oh",
        place_type=PlaceType.city.value,
        country_id=str(default_country.id),
        state_id=str(default_state.id),
    )


async def test_update_unknown_place_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
    default_state,
):
    resp = await authorized_client.put(
        f"/world/places/{uuid.uuid4()}",
        json=dict(
            name="X",
            code="Y",
            place_type=PlaceType.city.value,
            country_id=str(default_country.id),
            state_id=str(default_state.id),
        ),
    )
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "data": None,
        "message": "Place not found",
        "status_code": 404,
    }


async def test_update_place_with_unknown_country_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
    default_state,
):
    resp = await authorized_client.put(
        f"/world/places/{uuid.uuid4()}",
        json=dict(
            name="X",
            code="Y",
            place_type=PlaceType.city.value,
            country_id=str(uuid.uuid4()),
            state_id=str(default_state.id),
        ),
    )
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "data": None,
        "message": "Country not found",
        "status_code": 404,
    }


async def test_update_place_with_unknown_state_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
    default_state,
):
    resp = await authorized_client.put(
        f"/world/places/{uuid.uuid4()}",
        json=dict(
            name="X",
            code="Y",
            place_type=PlaceType.city.value,
            country_id=str(default_country.id),
            state_id=str(uuid.uuid4()),
        ),
    )
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "data": None,
        "message": "State not found",
        "status_code": 404,
    }


async def test_delete_place_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
    default_country,
    default_state,
):
    places = []
    for n in range(10):
        resp = await authorized_client.post(
            "/world/places",
            json=dict(
                name=f"Kacy Saj {n}",
                code=f"ks{n}",
                place_type=PlaceType.village.value,
                country_id=str(default_country.id),
                state_id=str(default_state.id),
            ),
        )
        assert resp.status_code == status.HTTP_201_CREATED
        places.append(resp.json())

    assert len(places) == 10
    await authorized_client.delete(f"/world/places/{places[0]['id']}")
    await authorized_client.delete(f"/world/places/{places[1]['id']}")

    resp = await authorized_client.get("/world/places")
    assert resp.status_code == status.HTTP_200_OK
    assert len(resp.json()) == 8
