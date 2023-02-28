FROM python:3.11-slim as builder

WORKDIR /build

COPY pyproject.toml poetry.lock poetry.toml ./
COPY alembic.ini alembic.ini
COPY migrations/ ./migrations
COPY observer/ ./observer

RUN pip install --no-cache-dir -U pip poetry
RUN poetry export --without-hashes --format=requirements.txt > requirements.txt

FROM builder

EXPOSE 8000
WORKDIR /observer

COPY --from=builder /build /observer
RUN pip install --no-cache-dir -r requirements.txt

CMD ["python", "-m", "observer", "server", "start", "--port", "8000"]
