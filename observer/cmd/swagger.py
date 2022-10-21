import json
from pathlib import Path

from fastapi import FastAPI
from typer import Option, Typer

from observer.api.health import router as health
from observer.settings import settings

swagger = Typer()


@swagger.command()
def generate(
    output_file: Path = Option(settings.swagger_output_file, help="Output file to save OpenAPI spec"),
):
    app = FastAPI(
        debug=settings.debug,
        title=settings.title,
        description=settings.description,
        redoc_url="/docs",
    )
    app.include_router(health)

    if not output_file.parent.exists():
        output_file.parent.mkdir(exist_ok=True)

    spec = app.openapi()
    with open(output_file, "w") as fp:
        fp.write(json.dumps(spec))
