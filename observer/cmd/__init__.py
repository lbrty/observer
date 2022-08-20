from typer import Typer

from observer.cmd.db import db
from observer.cmd.keys import keys

cli = Typer()

cli.add_typer(db, name="db")
cli.add_typer(keys, name="keys")
