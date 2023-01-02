import pytest
from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key
from cryptography.hazmat.primitives.serialization import (
    Encoding,
    NoEncryption,
    PrivateFormat,
)

from observer.services.keys import Keychain
from observer.services.storage import S3Storage


@pytest.mark.moto
async def test_keychain_can_load_keys_from_remote_store(aws_credentials, s3_server, s3_client, env_settings):
    bucket_name = "test-buck"
    await s3_client.create_bucket(
        Bucket=bucket_name,
        CreateBucketConfiguration=dict(LocationConstraint="eu-central-1"),
    )
    storage = S3Storage(bucket_name, env_settings.s3_region, s3_server)
    keychain = Keychain(storage)

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

        await s3_client.put_object(Bucket=bucket_name, Key=f"keys/key{n}.pem", Body=private_key_bytes)

    await keychain.load("keys")
    assert len(keychain.keys) == 2
