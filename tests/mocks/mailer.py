from typing import List

from observer.services.mailer import EmailMessage, IMailer


class MockMailer(IMailer):
    messages: List[EmailMessage] = []

    async def send(self, message: EmailMessage):
        self.messages.append(message)
