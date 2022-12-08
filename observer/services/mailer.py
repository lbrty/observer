from dataclasses import dataclass
from typing import Protocol


@dataclass
class Message:
    to_email: str
    from_email: str
    subject: str
    body: str


class MailerInterface(Protocol):
    async def send(self, message: Message):
        raise NotImplementedError
