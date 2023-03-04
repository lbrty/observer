image := sultaniman/observer
vsn := $(shell git log -1 --pretty=%h)

.PHONY: docker-image
docker-image:
	docker build . -t $(image):latest-slim -t $(image):$(vsn)-slim
	docker build . -f Dockerfile.alpine -t $(image):latest-alpine -t $(image):$(vsn)-alpine

.PHONY: fmt
fmt:
	poetry run black .
	poetry run isort .
	poetry run autoflake .

.PHONY: lint
lint:
	poetry run ruff . -q
	poetry run black . --check -q
	poetry run isort . -c -q
	poetry run mypy . --pretty

.PHONY: test
test:
	poetry run pytest tests --durations=0

.PHONY: serve
serve:
	poetry run python -m observer server start

.PHONY: swagger
swagger:
	poetry run python -m observer swagger generate
