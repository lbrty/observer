from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key
from cryptography.hazmat.primitives.serialization import BestAvailableEncryption, Encoding, PrivateFormat
from typer import Typer, Option

from observer.settings import settings

keys = Typer()


@keys.command()
def generate(
    filename: str = Option(..., "--output", "-o", help="Output file"),
    key_size: int = Option(settings.key_size, "--size", "-s", help="Size of RSA key"),
    password: str = Option(..., "--password", "-p", help="Password for key"),
):
    print(f"Generating key {filename}")

    priv_key = generate_private_key(
        public_exponent=settings.public_exponent,
        key_size=key_size,
    )

    priv_key_bytes = priv_key.private_bytes(
        encoding=Encoding.PEM,
        format=PrivateFormat.PKCS8,
        encryption_algorithm=BestAvailableEncryption(bytes(password)) if password else None
    )

    with open(filename, "wb") as fp:
        fp.write(priv_key_bytes)

    print("Done")
