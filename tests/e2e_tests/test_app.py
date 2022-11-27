from starlette import status

from observer.context import ctx


async def test_app(client):
    response = await client.get("/health")
    assert response.status_code == status.HTTP_200_OK
    assert response.json() == {"status": "ok"}


async def test_app_keys_loaded(app_context):
    assert len(ctx.key_loader.keys) == 1
