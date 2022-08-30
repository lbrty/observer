from alembic import command
from alembic.config import Config
from typer import Option, Typer

from observer.settings import db_settings, settings

db = Typer()


def load_config(uri) -> Config:
    config = Config(str(settings.base_path / "alembic.ini"))
    config.set_main_option("sqlalchemy.url", uri)
    config.set_main_option("script_location", str(settings.base_path / "migrations"))
    return config


@db.command()
def upgrade(
    uri: str = Option(db_settings.db_uri, help="Database URI DSN"),
    rev: str = Option("head", help="Revision to upgrade"),
):
    config = load_config(uri)
    command.upgrade(config, revision=rev)


@db.command()
def downgrade(
    uri: str = Option(db_settings.db_uri, help="Database URI DSN"),
    rev: str = Option("head", help="Revision to downgrade"),
):
    config = load_config(uri)
    command.downgrade(config, revision=rev)


@db.command()
def revision(
    uri: str = Option(db_settings.db_uri, help="Database URI DSN"),
    message: str = Option(..., "--message", "-m", help="Migration message"),
    auto: bool = Option(False, "--auto", "-a", is_flag=True, help="Auto generate migration"),
):
    config = load_config(uri)
    command.revision(config, message=message, autogenerate=auto)


@db.command()
def current(
    uri: str = Option(db_settings.db_uri, help="Database URI DSN"),
    verbose: bool = Option(False, is_flag=True, help="Verbose logs"),
):
    config = load_config(uri)
    command.current(config, verbose=verbose)


@db.command()
def history(
    uri: str = Option(db_settings.db_uri, help="Database URI DSN"),
    verbose: bool = Option(False, is_flag=True, help="Verbose logs"),
):
    config = load_config(uri)
    command.history(config, verbose=verbose, indicate_current=True)
