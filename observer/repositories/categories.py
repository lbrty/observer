from typing import List, Optional, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.people import categories
from observer.entities.people import Category, NewCategory, UpdateCategory


class ICategoryRepository(Protocol):
    async def create_category(self, new_category: NewCategory) -> Category:
        raise NotImplementedError

    async def get_categories(self, name: Optional[str] = None) -> List[Category]:
        raise NotImplementedError

    async def get_category(self, category_id: Identifier) -> Optional[Category]:
        raise NotImplementedError

    async def update_category(self, category_id: Identifier, updates: UpdateCategory) -> Optional[Category]:
        raise NotImplementedError

    async def delete_category(self, category_id: Identifier) -> Optional[Category]:
        raise NotImplementedError


class CategoryRepository(ICategoryRepository):
    def __init__(self, db: Database):
        self.db = db

    # Categories
    async def create_category(self, new_category: NewCategory) -> Category:
        # TODO: Handle IntegrityError for unique category names
        query = insert(categories).values(**new_category.dict()).returning("*")
        result = await self.db.fetchone(query)
        return Category(**result)

    async def get_categories(self, name: Optional[str] = None) -> List[Category]:
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
