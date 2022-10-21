from typer import Typer

from observer.cmd.db import db
from observer.cmd.keys import keys
from observer.cmd.server import server
from observer.cmd.swagger import swagger

cli = Typer()

cli.add_typer(db, name="db")
cli.add_typer(keys, name="keys")
cli.add_typer(server, name="server")
cli.add_typer(swagger, name="swagger")
