import { Button } from "@/components/ui/button";
import { type Column } from "@/components/table/data-table";
import { PencilSimpleIcon, UserCircleIcon } from "@/components/ui/icons";
import { StatusBadge } from "@/components/ui/status-badge";
import { TagChips } from "@/components/tags/tag-chips";
import type { Person } from "@/types/person";

interface PeopleColumnsOptions {
  projectId: string;
  t: (key: string) => string;
  onEdit: (id: string) => void;
  canWrite: boolean;
  statusLabels: Record<string, string>;
}

export function buildPeopleColumns({
  projectId,
  t,
  onEdit,
  canWrite,
  statusLabels,
}: PeopleColumnsOptions): Column<Person>[] {
  const columns: Column<Person>[] = [
    {
      key: "name",
      header: t("project.people.name"),
      render: (p) => (
        <div className="flex items-center gap-3">
          <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary text-fg-tertiary">
            <UserCircleIcon size={16} />
          </span>
          <div className="min-w-0">
            <p className="truncate font-medium text-fg">
              {p.first_name}
              {p.last_name ? ` ${p.last_name}` : ""}
            </p>
          </div>
        </div>
      ),
    },
    {
      key: "sex",
      header: t("project.people.sex"),
      render: (p) => <span className="text-fg-secondary">{p.sex}</span>,
    },
    {
      key: "case_status",
      header: t("project.people.caseStatus"),
      render: (p) => {
        const caseVariants: Record<string, "foam" | "gold" | "rose" | "neutral"> = {
          new: "gold",
          active: "foam",
          closed: "rose",
          archived: "neutral",
        };
        return (
          <StatusBadge
            label={statusLabels[p.case_status] ?? p.case_status}
            variant={caseVariants[p.case_status]}
          />
        );
      },
    },
    {
      key: "tags",
      header: t("project.tags.title"),
      render: (p) => <TagChips projectId={projectId} tagIds={p.tag_ids} />,
    },
    {
      key: "registered",
      header: t("project.people.registered"),
      render: (p) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(p.registered_at ?? p.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
  ];

  if (canWrite) {
    columns.push({
      key: "actions",
      header: "",
      render: (p: Person) => (
        <Button
          variant="ghost"
          className="p-1.5"
          onClick={(e) => {
            e.stopPropagation();
            onEdit(p.id);
          }}
        >
          <PencilSimpleIcon size={16} />
        </Button>
      ),
    } satisfies Column<Person>);
  }

  return columns;
}
