FROM python:3.11-slim

EXPOSE 8000
WORKDIR /observer

COPY pyproject.toml poetry.lock poetry.toml ./

RUN pip install -U pip poetry
RUN poetry export --without-hashes --format=requirements.txt > requirements.txt
RUN pip install -r requirements.txt

COPY alembic.ini ./alembic.ini
COPY migrations/ ./migrations
COPY observer/ ./observer

CMD ["python", "-m", "observer", "server", "start", "--port", "8000"]
