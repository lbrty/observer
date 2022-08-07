from typer import Typer

from observer.cmd.db import db

cli = Typer()

cli.add_typer(db, name="db")
