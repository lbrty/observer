from typing import List, Optional, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.people import Category, NewCategory, UpdateCategory
from observer.repositories.categories import ICategoryRepository
from observer.schemas.people import NewCategoryRequest, UpdateCategoryRequest


class ICategoryService(Protocol):
    repo: ICategoryRepository

    # Categories
    async def create_category(self, new_category: NewCategoryRequest) -> Category:
        raise NotImplementedError

    async def get_categories(self, name: Optional[str] = None) -> List[Category]:
        raise NotImplementedError

    async def get_category(self, category_id: Identifier) -> Category:
        raise NotImplementedError

    async def update_category(self, category_id: Identifier, updates: UpdateCategoryRequest) -> Category:
        raise NotImplementedError

    async def delete_category(self, category_id: Identifier) -> Category:
        raise NotImplementedError


class CategoryService(ICategoryService):
    def __init__(self, categories_repository: ICategoryRepository):
        self.repo = categories_repository

    # Categories
    async def create_category(self, new_category: NewCategoryRequest) -> Category:
        return await self.repo.create_category(NewCategory(**new_category.dict()))

    async def get_categories(self, name: Optional[str] = None) -> List[Category]:
        return await self.repo.get_categories(name)

    async def get_category(self, category_id: Identifier) -> Category:
        if category := await self.repo.get_category(category_id):
            return category

        raise NotFoundError(message="Category not found")

    async def update_category(self, category_id: Identifier, updates: UpdateCategoryRequest) -> Category:
        category = await self.repo.update_category(category_id, UpdateCategory(**updates.dict()))
        if category:
            return category

        raise NotFoundError(message="Category not found")

    async def delete_category(self, category_id: Identifier) -> Category:
        if category := await self.repo.delete_category(category_id):
            return category

        raise NotFoundError(message="Category not found")
