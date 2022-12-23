from typing import Protocol

from observer.repositories.idp import IDPRepositoryInterface


class IDPServiceInterface(Protocol):
    tag: str
    repo: IDPRepositoryInterface


class IDPService(IDPServiceInterface):
    tag: str = "source=service:idp"

    def __init__(self, idp_repository: IDPRepositoryInterface):
        self.repo = idp_repository
