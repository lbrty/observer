from typing import Any, Literal, Sequence, Type, get_args


def unwrap_literal(lit: Literal) -> Sequence[Any]:
    return get_args(lit)
