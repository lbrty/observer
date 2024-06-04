from sqlalchemy import func


utc_now = func.timezone("UTC", func.current_timestamp())
