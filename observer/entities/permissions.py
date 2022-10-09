from pydantic import BaseModel

from observer.common.types import Identifier


class Document(BaseModel):
    id: Identifier
    can_create: bool
    can_read: bool
    can_update: bool
    can_delete: bool
    can_create_projects: bool
    can_read_documents: bool
    can_read_personal_info: bool
    user_id: Identifier
    project_id: Identifier
