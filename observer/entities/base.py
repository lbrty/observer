from pydantic import BaseModel


class ModelBase(BaseModel):
    class Config:
        frozen = True
        use_enum_values = True
