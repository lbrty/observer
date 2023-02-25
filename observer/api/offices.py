from typing import Optional

from fastapi import APIRouter, BackgroundTasks, Depends, Query, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.types import Identifier, Role
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.pagination import pagination
from observer.components.services import audit_service, office_service
from observer.entities.users import User
from observer.schemas.offices import (
    NewOfficeRequest,
    OfficeResponse,
    OfficesResponse,
    UpdateOfficeRequest,
)
from observer.schemas.pagination import Pagination
from observer.services.audit_logs import IAuditService
from observer.services.offices import IOfficesService

router = APIRouter(prefix="/offices")


@router.post(
    "",
    response_model=OfficeResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
    ),
    tags=["offices"],
)
async def create_office(
    tasks: BackgroundTasks,
    new_office: NewOfficeRequest,
    user: Optional[User] = Depends(
        RequiresRoles([Role.admin, Role.staff]),
    ),
    offices: IOfficesService = Depends(office_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_office,action=create:office",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> OfficeResponse:
    office = await offices.create_office(new_office.name)
    audit_log = props.new_event(f"office_id={office.id},ref_id={user.id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return OfficeResponse(**office.dict())


@router.get(
    "",
    response_model=OfficesResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.staff, Role.consultant]),
        ),
    ],
    tags=["offices"],
)
async def get_offices(
    name: Optional[str] = Query(None, description="When given offices will be filtered by name"),
    offices: IOfficesService = Depends(office_service),
    pages: Pagination = Depends(pagination),
) -> OfficesResponse:
    total, items = await offices.get_offices(name, pages.offset, pages.limit)
    return OfficesResponse(total=total, items=items)


@router.get(
    "/{office_id}",
    response_model=OfficeResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.staff, Role.consultant]),
        ),
    ],
    tags=["offices"],
)
async def get_office(
    office_id: Identifier,
    offices: IOfficesService = Depends(office_service),
) -> OfficeResponse:
    office = await offices.get_office(office_id)
    return OfficeResponse(**office.dict())


@router.put(
    "/{office_id}",
    response_model=OfficeResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.staff]),
        ),
    ],
    tags=["offices"],
)
async def update_office(
    tasks: BackgroundTasks,
    office_id: Identifier,
    updates: UpdateOfficeRequest,
    user: Optional[User] = Depends(
        RequiresRoles([Role.admin, Role.staff]),
    ),
    offices: IOfficesService = Depends(office_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_office,action=update:office",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> OfficeResponse:
    office = await offices.get_office(office_id)
    if office.name != updates.name:
        office = await offices.update_office(office_id, updates.name)
        audit_log = props.new_event(
            f"office_id={office.id},ref_id={user.id}",
            jsonable_encoder(office),
        )
        tasks.add_task(audits.add_event, audit_log)
    return OfficeResponse(**office.dict())


@router.delete(
    "/{office_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.staff]),
        ),
    ],
    tags=["offices"],
)
async def delete_office(
    tasks: BackgroundTasks,
    office_id: Identifier,
    user: Optional[User] = Depends(
        RequiresRoles([Role.admin, Role.staff]),
    ),
    offices: IOfficesService = Depends(office_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_office,action=delete:office",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    office = await offices.delete_office(office_id)
    audit_log = props.new_event(f"office_id={office.id},ref_id={user.id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
