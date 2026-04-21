import type { Column } from "@/components/data-table";
import { StatusBadge, StatusDot } from "@/components/status-badge";
import { UserInitials } from "@/components/user-initials";
import type { AdminUser } from "@/types/admin";

interface UsersColumnsOptions {
  t: (key: string) => string;
}

export function buildUsersColumns({ t }: UsersColumnsOptions): Column<AdminUser>[] {
  return [
    {
      key: "name",
      header: t("admin.users.name"),
      render: (u) => (
        <div className="flex items-center gap-3">
          <UserInitials firstName={u.first_name} lastName={u.last_name} />
          <span className="font-medium text-fg">
            {u.first_name} {u.last_name}
          </span>
        </div>
      ),
    },
    {
      key: "email",
      header: t("admin.users.email"),
      render: (u) => <span className="text-fg-secondary">{u.email}</span>,
    },
    {
      key: "role",
      header: t("admin.users.role"),
      render: (u) => <StatusBadge label={u.role} />,
    },
    {
      key: "active",
      header: t("admin.users.active"),
      render: (u) =>
        u.deactivated_at ? (
          <span className="inline-flex items-center gap-1 text-xs text-amber-600">
            <span className="size-1.5 rounded-full bg-amber-500" />
            {t("users.deactivated")}
          </span>
        ) : (
          <StatusDot active={u.is_active} />
        ),
    },
    {
      key: "verified",
      header: t("admin.users.verified"),
      render: (u) => <StatusDot active={u.is_verified} />,
    },
    {
      key: "created",
      header: t("admin.users.created"),
      render: (u) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(u.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
  ];
}
