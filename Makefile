image := sultaniman/observer
vsn := $(shell git log -1 --pretty=%h)

.PHONY: docker-image
docker-image:
	docker build . -t $(image):latest-slim -t $(image):$(vsn)-slim
	docker build . -f Dockerfile.alpine -t $(image):latest-alpine -t $(image):$(vsn)-alpine

.PHONY: fmt
fmt:
	ruff format .

.PHONY: lint
lint:
	ruff . -q
	mypy . --pretty

.PHONY: test
test:
	pytest tests

.PHONY: serve
serve:
	python -m observer server start

.PHONY: swagger
swagger:
	python -m observer swagger generate
