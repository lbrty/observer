from typing import List, Optional, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import Identifier, SomeStr
from observer.db import Database
from observer.db.tables.idp import categories
from observer.entities.idp import Category, NewCategory, UpdateCategory


class CategoryRepositoryInterface(Protocol):
    async def create_category(self, new_category: NewCategory) -> Category:
        raise NotImplementedError

    async def get_categories(self, name: SomeStr = None) -> List[Category]:
        raise NotImplementedError

    async def get_category(self, category_id: Identifier) -> Optional[Category]:
        raise NotImplementedError

    async def update_category(self, category_id: Identifier, updates: UpdateCategory) -> Optional[Category]:
        raise NotImplementedError

    async def delete_category(self, category_id: Identifier) -> Optional[Category]:
        raise NotImplementedError


class CategoryRepository(CategoryRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    # Categories
    async def create_category(self, new_category: NewCategory) -> Category:
        # TODO: Handle IntegrityError for unique category names
        query = insert(categories).values(**new_category.dict()).returning("*")
        result = await self.db.fetchone(query)
        return Category(**result)

    async def get_categories(self, name: SomeStr = None) -> List[Category]:
        query = select(categories)
        if name:
            query = query.where(categories.c.name.ilike(f"%{name}%"))

        rows = await self.db.fetchall(query)
        return [Category(**row) for row in rows]

    async def get_category(self, category_id: Identifier) -> Optional[Category]:
        query = select(categories).where(categories.c.id == category_id)
        if row := await self.db.fetchone(query):
            return Category(**row)

        return None

    async def update_category(self, category_id: Identifier, updates: UpdateCategory) -> Optional[Category]:
        query = update(categories).values(**updates.dict()).where(categories.c.id == category_id).returning("*")
        if row := await self.db.fetchone(query):
            return Category(**row)

        return None

    async def delete_category(self, category_id: Identifier) -> Optional[Category]:
        query = delete(categories).where(categories.c.id == category_id).returning("*")
        if row := await self.db.fetchone(query):
            return Category(**row)

        return None
