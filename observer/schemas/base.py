from pydantic import BaseModel


class SchemaBase(BaseModel):
    class Config:
        use_enum_values = True
