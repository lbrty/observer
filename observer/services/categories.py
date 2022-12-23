from datetime import datetime, timedelta, timezone
from typing import List, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier, SomeStr
from observer.entities.idp import Category, NewCategory, UpdateCategory
from observer.repositories.categories import CategoryRepositoryInterface
from observer.schemas.audit_logs import NewAuditLog
from observer.schemas.idp import NewCategoryRequest, UpdateCategoryRequest


class CategoryServiceInterface(Protocol):
    tag: str
    repo: CategoryRepositoryInterface

    # Categories
    async def create_category(self, new_category: NewCategoryRequest) -> Category:
        raise NotImplementedError

    async def get_categories(self, name: SomeStr = None) -> List[Category]:
        raise NotImplementedError

    async def get_category(self, category_id: Identifier) -> Category:
        raise NotImplementedError

    async def update_category(self, category_id: Identifier, updates: UpdateCategoryRequest) -> Category:
        raise NotImplementedError

    async def delete_category(self, category_id: Identifier) -> Category:
        raise NotImplementedError

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        raise NotImplementedError


class CategoryService(CategoryServiceInterface):
    tag: str = "source=service:categories"

    def __init__(self, categories_repository: CategoryRepositoryInterface):
        self.repo = categories_repository

    # Categories
    async def create_category(self, new_category: NewCategoryRequest) -> Category:
        return await self.repo.create_category(NewCategory(**new_category.dict()))

    async def get_categories(self, name: SomeStr = None) -> List[Category]:
        return await self.repo.get_categories(name)

    async def get_category(self, category_id: Identifier) -> Category:
        if category := await self.repo.get_category(category_id):
            return category

        raise NotFoundError(message="Category not found")

    async def update_category(self, category_id: Identifier, updates: UpdateCategoryRequest) -> Category:
        await self.get_category(category_id)
        return await self.repo.update_category(category_id, UpdateCategory(**updates.dict()))

    async def delete_category(self, category_id: Identifier) -> Category:
        await self.get_category(category_id)
        return await self.repo.delete_category(category_id)

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        now = datetime.now(tz=timezone.utc)
        expires_at = None
        if expires_in:
            expires_at = now + expires_in

        return NewAuditLog(
            ref=f"{self.tag},{ref}",
            data=data,
            expires_at=expires_at,
        )
