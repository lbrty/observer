from dataclasses import dataclass
from typing import Protocol


@dataclass
class EmailMessage:
    to_email: str
    from_email: str
    subject: str
    body: str


class MailerInterface(Protocol):
    async def send(self, message: EmailMessage):
        raise NotImplementedError


class Mailer(MailerInterface):
    async def send(self, message: EmailMessage):
        pass
