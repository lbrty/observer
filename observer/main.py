from observer.app import create_app
from observer.settings import settings

app = create_app(settings)


@app.on_event("startup")
def on_startup():
    pass


@app.on_event("shutdown")
def on_shutdown():
    pass
