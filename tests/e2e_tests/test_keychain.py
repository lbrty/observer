from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key
from cryptography.hazmat.primitives.serialization import (
    Encoding,
    NoEncryption,
    PrivateFormat,
)

from observer.services.keychain import Keychain
from observer.services.storage import FSStorage


async def test_keychain_can_load_keys_from_remote_store(
    aws_credentials,
    create_bucket,
    s3_storage,
    s3_client,
    app_context,
    env_settings,
):
    bucket_name = "test-buck"
    keychain = Keychain()

    for n in range(2):
        private_key = generate_private_key(
            public_exponent=env_settings.public_exponent,
            key_size=env_settings.key_size,
        )

        private_key_bytes = private_key.private_bytes(
            encoding=Encoding.PEM,
            format=PrivateFormat.PKCS8,
            encryption_algorithm=NoEncryption(),
        )

        await s3_client.put_object(Bucket=bucket_name, Key=f"uploads/keys/key{n}.pem", Body=private_key_bytes)

    await keychain.load("keys", app_context.storage)
    assert len(keychain.keys) == 2


async def test_keychain_can_load_keys_from_filesystem_store(temp_keystore, env_settings):
    env_settings.keystore_path = temp_keystore
    storage = FSStorage(env_settings.keystore_path)
    keychain = Keychain()
    await keychain.load(temp_keystore, storage)
    assert len(keychain.keys) == 5
