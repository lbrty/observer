from alembic import command
from alembic.config import Config
from typer import Typer, Option

from observer.settings import settings

db = Typer()

config = Config(settings.base_path / "alembic.ini")
config.set_main_option("sqlalchemy.url", settings.db_uri)


@db.command()
def upgrade(
    uri: str = Option(settings.db_uri, help="Database URI DSN"),
    rev: str = Option("head", help="Revision to upgrade"),
):
    config.set_main_option("sqlalchemy.url", uri)
    command.upgrade(config, revision=rev)


@db.command()
def downgrade(
    uri: str = Option(settings.db_uri, help="Database URI DSN"),
    rev: str = Option("head", help="Revision to downgrade"),
):
    config.set_main_option("sqlalchemy.url", uri)
    command.downgrade(config, revision=rev)


@db.command()
def revision(
    uri: str = Option(settings.db_uri, help="Database URI DSN"),
    message: str = Option(..., "--message", "-m", help="Migration message"),
    auto: bool = Option(False, "--auto", "-a", is_flag=True, help="Auto generate migration"),
):
    config.set_main_option("sqlalchemy.url", uri)
    command.revision(config, message=message, autogenerate=auto)


@db.command()
def current(
    uri: str = Option(settings.db_uri, help="Database URI DSN"),
    verbose: bool = Option(False, is_flag=True, help="Verbose logs"),
):
    config.set_main_option("sqlalchemy.url", uri)
    command.current(config, verbose=verbose)


@db.command()
def history(
    uri: str = Option(settings.db_uri, help="Database URI DSN"),
    verbose: bool = Option(False, is_flag=True, help="Verbose logs"),
):
    config.set_main_option("sqlalchemy.url", uri)
    command.history(config, verbose=verbose, indicate_current=True)
