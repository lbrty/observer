from typing import Optional

from sqlalchemy import Table, asc, desc
from sqlalchemy.orm import Query


def parse_order_by(field: str, table: Table) -> Optional[Query]:
    real_field = field
    if field.startswith("-"):
        real_field = field[1:]
        fun = desc
    else:
        fun = asc
        if field.startswith("+"):
            real_field = field[1:]

    if real_field in table.c:
        return fun(table.c[real_field])

    return None
