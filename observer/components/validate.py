from fastapi import Header

from observer.api.exceptions import TooLargeDocumentError


class ContentLengthLimit:
    def __init__(self, max_size: int):
        self.max_size = max_size

    async def __call__(
        self,
        content_length: int = Header(..., description="Uploaded file size"),
    ):
        if content_length > self.max_size:
            raise TooLargeDocumentError
