import datetime

today = datetime.date.today()

project = "Observer"
copyright = f"{today.year}, Sultan Iman"  # noqa: A001
author = "Sultan Iman"
release = "0.5.0"

extensions = [
    "sphinxawesome_theme.highlighting",
    "sphinx.ext.napoleon",
]

templates_path = ["_templates"]
exclude_patterns = ["_build", "Thumbs.db", ".DS_Store"]

html_theme = "sphinxawesome_theme"
html_static_path = ["_static"]
