image := sultaniman/observer
vsn := $(shell git log -1 --pretty=%h)

.PHONY: docker-image
docker-image:
	docker build . -t $(image):latest-slim -t $(image):$(vsn)-slim
	docker build . -f Dockerfile.alpine -t $(image):latest-alpine -t $(image):$(vsn)-alpine

.PHONY: fmt
fmt:
	poetry run ruff format .

.PHONY: lint
lint:
	poetry run ruff . -q
	poetry run mypy . --pretty

.PHONY: test
test:
	poetry run pytest tests

.PHONY: serve
serve:
	poetry run python -m observer server start

.PHONY: swagger
swagger:
	poetry run python -m observer swagger generate
