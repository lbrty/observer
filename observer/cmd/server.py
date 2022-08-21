import uvicorn
from typer import Option, Typer

from observer.settings import settings

server = Typer()


@server.command()
def start(port: int = Option(settings.port, help="Port to bind server")):
    """Start API server"""
    uvicorn.run(
        "observer.main:app",
        host="0.0.0.0",
        port=port,
        debug=settings.debug,
        reload=settings.debug,
    )
