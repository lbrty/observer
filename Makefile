.PHONY: fmt
fmt:
	poetry run black .
	poetry run isort .
	poetry run autoflake .

.PHONY: lint
lint:
	#TODO: bug in github poetry run autoflake . -c --quiet -j 10
	poetry run black . --check --quiet
	poetry run isort . -c --quiet

.PHONY: serve
serve:
	python -m observer server start
