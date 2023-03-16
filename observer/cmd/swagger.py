import json
from pathlib import Path

from typer import Option, Typer

from observer.app import create_app
from observer.settings import settings

swagger = Typer()


@swagger.command()
def generate(
    output_file: Path = Option(settings.swagger_output_file, help="Output file to save OpenAPI spec"),
):
    """Generate OpenAPI spec"""
    app = create_app(settings, None)
    if not output_file.parent.exists():
        output_file.parent.mkdir(exist_ok=True)

    spec = app.openapi()
    with open(output_file, "w") as fp:
        fp.write(json.dumps(spec, indent=2))
