import os
from dataclasses import dataclass
from typing import Protocol
from urllib.error import HTTPError

from sendgrid import SendGridAPIClient
from sendgrid.helpers.mail import Content, From, Mail, To
from structlog import get_logger

from observer.api.exceptions import InternalError

logger = get_logger("mailer")


@dataclass
class EmailMessage:
    to_email: str
    from_email: str
    subject: str
    body: str


class IMailer(Protocol):
    async def send(self, message: EmailMessage):
        raise NotImplementedError


class Mailer(IMailer):
    async def send(self, message: EmailMessage):
        pass


class SendgridMailer(IMailer):
    def __init__(self):
        if api_key := os.environ.get("SENDGRID_API_KEY"):
            self.client = SendGridAPIClient(api_key=api_key)
        else:
            raise InternalError(message="SENDGRID_API_KEY is not configured")

    async def send(self, message: EmailMessage):
        try:
            from_email = From(message.from_email)
            to_email = To(message.to_email)
            plain_text_content = Content("text/plain", message.body)

            # instantiate `sendgrid.helpers.mail.Mail` objects
            mail = Mail(from_email, to_email, plain_text_content)
            self.client.send(mail)
        except HTTPError as e:
            logger.exception("Unable send Sendgrid email", response=e.read().decode())
            raise InternalError(message="Unable send Sendgrid email")
