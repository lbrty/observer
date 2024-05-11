from observer.db.models import AuditLog
from observer.entities.audit_logs import NewAuditLog


class IAuditLogs:
    pass


class AuditLogs(IAuditLogs):
    """Repo to manage audit logs."""
    def add(self, event: NewAuditLog) -> AuditLog:
        pass
