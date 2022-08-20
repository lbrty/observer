import hashlib
from glob import glob
from typing import Optional

from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key
from cryptography.hazmat.primitives.serialization import (
    BestAvailableEncryption,
    Encoding,
    NoEncryption,
    PrivateFormat,
    load_pem_private_key,
)
from typer import Typer, Option

from observer.settings import settings

keys = Typer()


@keys.command()
def generate(
    filename: str = Option(..., "--output", "-o", help="Output file (should be filename i.e. key.pem)"),
    key_size: int = Option(settings.key_size, "--size", "-s", help="Size of RSA key"),
    password: Optional[str] = Option(None, "--password", "-p", help="Password for key"),
):
    print(f"Generating key {filename}")

    priv_key = generate_private_key(
        public_exponent=settings.public_exponent,
        key_size=key_size,
    )

    priv_key_bytes = priv_key.private_bytes(
        encoding=Encoding.PEM,
        format=PrivateFormat.PKCS8,
        encryption_algorithm=BestAvailableEncryption(bytes(password)) if password else NoEncryption(),
    )

    with open(settings.key_store_path / filename, "wb") as fp:
        fp.write(priv_key_bytes)

    print("Done")


@keys.command()
def list_keys():
    key_list = glob(f"{str(settings.key_store_path)}/*.pem")
    for key in key_list:
        with open(key, "rb") as fp:
            pem = load_pem_private_key(fp.read(), password=None)
            key_bytes = pem.private_bytes(
                encoding=Encoding.PEM,
                format=PrivateFormat.PKCS8,
                encryption_algorithm=NoEncryption(),
            )
            h = hashlib.new("sha256", key_bytes)
            print(f"{str(h.hexdigest())}:{key}")
