.PHONY: fmt
fmt:
	black . --line-length 120

.PHONY: serve
serve:
	python -m observer server start
