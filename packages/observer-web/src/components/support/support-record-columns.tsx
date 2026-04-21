import type { Column } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { HandHeartIcon, PencilSimpleIcon } from "@/components/ui/icons";
import { StatusBadge } from "@/components/ui/status-badge";
import { referralKeys, sphereKeys, typeKeys } from "@/constants/support";
import type { SupportRecord } from "@/types/support-record";

interface SupportRecordColumnsOptions {
  t: (key: string) => string;
  canWrite: boolean;
  onEdit: (id: string) => void;
}

export function buildSupportRecordColumns({
  t,
  canWrite,
  onEdit,
}: SupportRecordColumnsOptions): Column<SupportRecord>[] {
  return [
    {
      key: "person_id",
      header: t("project.supportRecords.person"),
      render: (r) => (
        <div className="flex items-center gap-3">
          <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary text-fg-tertiary">
            <HandHeartIcon size={16} />
          </span>
          <span className="truncate text-sm text-fg">
            {[r.person_first_name, r.person_last_name].filter(Boolean).join(" ") || r.person_id}
          </span>
        </div>
      ),
    },
    {
      key: "type",
      header: t("project.supportRecords.type"),
      render: (r) => (
        <StatusBadge label={typeKeys[r.type] ? t(typeKeys[r.type]) : r.type} statusKey={r.type} />
      ),
    },
    {
      key: "sphere",
      header: t("project.supportRecords.sphere"),
      render: (r) => (
        <span className="text-fg-secondary">
          {r.sphere ? t(sphereKeys[r.sphere] ?? r.sphere) : "\u2014"}
        </span>
      ),
    },
    {
      key: "provided_at",
      header: t("project.supportRecords.providedAt"),
      render: (r) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {r.provided_at ? new Date(r.provided_at).toLocaleDateString("en-CA") : "\u2014"}
        </span>
      ),
    },
    {
      key: "referral_status",
      header: t("project.supportRecords.referralStatus"),
      render: (r) =>
        r.referral_status ? (
          <StatusBadge
            label={
              referralKeys[r.referral_status]
                ? t(referralKeys[r.referral_status])
                : r.referral_status
            }
            statusKey={r.referral_status}
          />
        ) : (
          <span className="text-fg-tertiary">{"\u2014"}</span>
        ),
    },
    ...(canWrite
      ? [
          {
            key: "actions",
            header: "",
            render: (r: SupportRecord) => (
              <Button
                variant="ghost"
                className="p-1.5"
                onClick={(e) => {
                  e.stopPropagation();
                  onEdit(r.id);
                }}
              >
                <PencilSimpleIcon size={16} />
              </Button>
            ),
          } satisfies Column<SupportRecord>,
        ]
      : []),
  ];
}
