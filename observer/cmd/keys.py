import asyncio
import logging
from pathlib import Path

import structlog
from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key
from cryptography.hazmat.primitives.serialization import (
    Encoding,
    NoEncryption,
    PrivateFormat,
)
from rich.console import Console
from rich.tree import Tree
from typer import Option, Typer

from observer.services.keychain import Keychain
from observer.services.storage import init_storage
from observer.settings import settings

keys = Typer()
console = Console()

structlog.configure(
    wrapper_class=structlog.make_filtering_bound_logger(logging.WARNING),
)


@keys.command()
def generate(
    filename: str = Option(..., "--output", "-o", help="Output file (should be filename i.e. key.pem)"),
    key_size: int = Option(settings.key_size, "--size", "-s", help="Size of RSA key"),
):
    """Generate RSA key"""
    print(f"Generating key {filename}")

    priv_key = generate_private_key(
        public_exponent=settings.public_exponent,
        key_size=key_size,
    )

    priv_key_bytes = priv_key.private_bytes(
        encoding=Encoding.PEM,
        format=PrivateFormat.PKCS8,
        encryption_algorithm=NoEncryption(),
    )

    with open(Path(settings.keystore_path) / filename, "wb") as fp:
        fp.write(priv_key_bytes)

    print("Done")


@keys.command()
def list_keys():
    """List all keys from key store"""
    storage = init_storage(settings.storage_kind, settings)
    keychain = Keychain()
    asyncio.get_event_loop().run_until_complete(keychain.load(settings.keystore_path, storage))

    tree = Tree(f"Key store: {settings.keystore_path}")
    if keychain.keys:
        for key in keychain.keys:
            tree.add(f"{Path(key.filename).name} \t [blue bold]{key.hash}[/]")
    else:
        tree.add("[bold red]No keys found[/]")
    console.print(tree)
