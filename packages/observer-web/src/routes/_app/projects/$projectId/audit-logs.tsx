import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { ClockCounterClockwiseIcon } from "@/components/icons";
import type { Column } from "@/components/data-table";
import { DataTablePage } from "@/components/data-table-page";
import type { FilterDef } from "@/components/filter-bar";
import { useProjectAuditLogs } from "@/hooks/use-audit-logs";
import type { AuditEntry } from "@/types/audit";

export const Route = createFileRoute("/_app/projects/$projectId/audit-logs")({
  component: ProjectAuditLogsPage,
});

const ACTION_VALUES = [
  "create",
  "update",
  "delete",
  "export",
];

const ENTITY_TYPE_VALUES = [
  "person",
  "support_record",
  "migration_record",
  "household",
  "note",
  "document",
  "tag",
  "pet",
  "permission",
];

function ProjectAuditLogsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();

  const [page, setPage] = useState(1);
  const [action, setAction] = useState("");
  const [entityType, setEntityType] = useState("");
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");

  const params = {
    page,
    per_page: 20,
    ...(action && { action }),
    ...(entityType && { entity_type: entityType }),
    ...(dateFrom && { date_from: dateFrom }),
    ...(dateTo && { date_to: dateTo }),
  };

  const { data, isLoading } = useProjectAuditLogs(projectId, params);

  const actionOptions = [
    { label: t("audit.allActions"), value: "" },
    ...ACTION_VALUES.map((v) => ({ label: v, value: v })),
  ];

  const entityOptions = [
    { label: t("audit.allEntities"), value: "" },
    ...ENTITY_TYPE_VALUES.map((v) => ({ label: v, value: v })),
  ];

  const columns: Column<AuditEntry>[] = [
    {
      key: "created_at",
      header: t("audit.timestamp"),
      render: (e) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(e.created_at).toLocaleString("en-CA", {
            dateStyle: "short",
            timeStyle: "medium",
          })}
        </span>
      ),
    },
    {
      key: "action",
      header: t("audit.action"),
      render: (e) => (
        <span className="rounded bg-bg-tertiary px-1.5 py-0.5 text-xs font-medium text-fg-secondary">
          {e.action}
        </span>
      ),
    },
    {
      key: "entity_type",
      header: t("audit.entityType"),
      render: (e) => (
        <span className="text-sm text-fg-secondary">{e.entity_type}</span>
      ),
    },
    {
      key: "summary",
      header: t("audit.summary"),
      render: (e) => (
        <span className="max-w-xs truncate text-sm text-fg">{e.summary}</span>
      ),
    },
    {
      key: "ip",
      header: t("audit.ip"),
      render: (e) => (
        <span className="font-mono text-xs text-fg-tertiary">{e.ip}</span>
      ),
    },
  ];

  const filters: FilterDef[] = [
    {
      type: "select",
      value: action,
      onValueChange: (v) => { setAction(v); setPage(1); },
      options: actionOptions,
      placeholder: t("audit.allActions"),
    },
    {
      type: "select",
      value: entityType,
      onValueChange: (v) => { setEntityType(v); setPage(1); },
      options: entityOptions,
      placeholder: t("audit.allEntities"),
    },
    {
      type: "date-range",
      fromValue: dateFrom,
      toValue: dateTo,
      onFromChange: (v) => { setDateFrom(v); setPage(1); },
      onToChange: (v) => { setDateTo(v); setPage(1); },
      fromPlaceholder: t("common.dateFrom"),
      toPlaceholder: t("common.dateTo"),
    },
  ];

  return (
    <DataTablePage
      title={t("audit.title")}
      columns={columns}
      data={data?.entries ?? []}
      keyExtractor={(e) => e.id}
      isLoading={isLoading}
      filters={filters}
      pagination={data ? { page: data.page, perPage: data.per_page, total: data.total, onChange: setPage } : undefined}
      emptyIcon={ClockCounterClockwiseIcon}
      emptyTitle={t("audit.emptyTitle")}
    />
  );
}
