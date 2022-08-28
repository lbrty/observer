import uvicorn
from typer import Option, Typer

from observer.settings import settings

server = Typer()


@server.command()
def start(
    port: int = Option(settings.port, help="Port to bind server"),
    host: str = Option("0.0.0.0", envvar="HOST", help="Host to bind"),
):
    """Start API server"""
    uvicorn.run(
        "observer.main:app",
        host=host,
        port=port,
        debug=settings.debug,
        reload=settings.debug,
    )
