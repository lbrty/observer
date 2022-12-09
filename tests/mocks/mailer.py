from typing import List

from observer.services.mailer import EmailMessage, MailerInterface


class MockMailer(MailerInterface):
    messages: List[EmailMessage] = []

    async def send(self, message: EmailMessage):
        self.messages.append(message)
