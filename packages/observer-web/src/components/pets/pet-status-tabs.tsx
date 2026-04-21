import { Tabs } from "@base-ui/react/tabs";
import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import type { PetStatus } from "@/routes/_app/projects/$projectId/pets/-pets-page";

const statusTabs: PetStatus[] = [
  "",
  "registered",
  "adopted",
  "owner_found",
  "needs_shelter",
  "unknown",
];

export function PetStatusTabs({
  projectId,
  statusFilter,
  statusLabels,
}: {
  projectId: string;
  statusFilter: PetStatus;
  statusLabels: Record<string, string>;
}) {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <Tabs.Root
      defaultValue=""
      value={statusFilter}
      onValueChange={(value) => {
        const s = value as PetStatus;
        if (s) {
          navigate({ to: "/projects/$projectId/pets/$status", params: { projectId, status: s } });
        } else {
          navigate({ to: "/projects/$projectId/pets", params: { projectId } });
        }
      }}
      className="mb-4"
    >
      <Tabs.List className="flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
        {statusTabs.map((tab) => (
          <Tabs.Tab
            key={tab}
            value={tab}
            className="cursor-pointer rounded-sm px-4 py-1.5 m-0.5 text-sm font-medium text-fg-tertiary transition-colors hover:text-fg data-active:bg-bg data-active:text-fg data-active:shadow-card"
          >
            {tab === "" ? t("project.pets.all") : statusLabels[tab]}
          </Tabs.Tab>
        ))}
      </Tabs.List>
    </Tabs.Root>
  );
}
