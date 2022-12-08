import hashlib
from glob import glob
from pathlib import Path

from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key
from cryptography.hazmat.primitives.serialization import (
    Encoding,
    NoEncryption,
    PrivateFormat,
)
from rich.console import Console
from rich.tree import Tree
from typer import Option, Typer

from observer.settings import settings

keys = Typer()
console = Console()


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

    with open(settings.keystore_path / filename, "wb") as fp:
        fp.write(priv_key_bytes)

    print("Done")


@keys.command()
def list_keys():
    """List all keys from key store"""
    key_list = glob(f"{str(settings.keystore_path)}/*.pem")
    tree = Tree(f"Key store: {settings.keystore_path}")
    if key_list:
        for key in key_list:
            with open(key, "rb") as fp:
                file_bytes = fp.read()
                h = hashlib.new("sha256", file_bytes)
                tree.add(f"{Path(key).name}    [blue bold]{str(h.hexdigest())[:16].upper()}[/]")
    else:
        tree.add("[bold red]No keys found[/]")
    console.print(tree)
