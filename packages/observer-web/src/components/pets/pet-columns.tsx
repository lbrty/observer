import { Button } from "@/components/ui/button";
import { DataTable, type Column } from "@/components/table/data-table";
import { PawPrintIcon, PencilSimpleIcon } from "@/components/ui/icons";
import { PersonName } from "@/components/people/person-name";
import { StatusBadge } from "@/components/ui/status-badge";
import { TagChips } from "@/components/tags/tag-chips";
import type { Pet } from "@/types/pet";

// Re-export Column so callers don't need to import data-table separately.
export type { Column };

interface BuildPetColumnsOptions {
  t: (key: string) => string;
  projectId: string;
  statusLabels: Record<string, string>;
  statusVariants: Record<string, "foam" | "gold" | "rose" | "neutral">;
  canWrite: boolean;
  onEdit: (id: string) => void;
}

export function buildPetColumns({
  t,
  projectId,
  statusLabels,
  statusVariants,
  canWrite,
  onEdit,
}: BuildPetColumnsOptions): Column<Pet>[] {
  const base: Column<Pet>[] = [
    {
      key: "name",
      header: t("project.pets.name"),
      render: (p) => (
        <div className="flex items-center gap-3">
          <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary text-fg-tertiary">
            <PawPrintIcon size={16} />
          </span>
          <div className="min-w-0">
            <p className="truncate font-medium text-fg">{p.name}</p>
          </div>
        </div>
      ),
    },
    {
      key: "status",
      header: t("project.pets.status"),
      render: (p) => (
        <StatusBadge
          label={statusLabels[p.status] ?? p.status}
          variant={statusVariants[p.status]}
        />
      ),
    },
    {
      key: "owner_id",
      header: t("project.pets.ownerId"),
      render: (p) =>
        p.owner_id ? (
          <span className="text-sm text-fg-secondary">
            <PersonName projectId={projectId} personId={p.owner_id} />
          </span>
        ) : (
          <span className="text-fg-tertiary">—</span>
        ),
    },
    {
      key: "tags",
      header: t("project.tags.title"),
      render: (p) => <TagChips projectId={projectId} tagIds={p.tag_ids} />,
    },
    {
      key: "registration_id",
      header: t("project.pets.registrationId"),
      render: (p) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {p.registration_id ?? ""}
        </span>
      ),
    },
  ];

  if (canWrite) {
    base.push({
      key: "actions",
      header: "",
      render: (p: Pet) => (
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
    } satisfies Column<Pet>);
  }

  return base;
}
