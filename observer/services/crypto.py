from typing import Any, Protocol, Sequence


class KeychainLoader(Protocol):
    keys: Sequence[Any] = []

    async def load(self, path: str):
        raise NotImplementedError

    @property
    def signatures(self) -> Sequence[str]:
        raise NotImplementedError
