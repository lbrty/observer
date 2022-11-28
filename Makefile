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

.PHONY: test
test:
	poetry run pytest tests

.PHONY: serve
serve:
	python -m observer server start

.PHONY: swagger
swagger:
	python -m observer swagger generate
