import { Tabs } from "@base-ui/react/tabs";
import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

const supportTypes = [
  "",
  "humanitarian",
  "legal",
  "social",
  "psychological",
  "medical",
  "general",
] as const;

interface SupportRecordStatusTabsProps {
  projectId: string;
  typeFilter: string;
}

export function SupportRecordStatusTabs({ projectId, typeFilter }: SupportRecordStatusTabsProps) {
  const { t } = useTranslation();

  const tabLabels: Record<string, string> = {
    "": t("project.supportRecords.all"),
    humanitarian: t("project.supportRecords.typeHumanitarian"),
    legal: t("project.supportRecords.typeLegal"),
    social: t("project.supportRecords.typeSocial"),
    psychological: t("project.supportRecords.typePsychological"),
    medical: t("project.supportRecords.typeMedical"),
    general: t("project.supportRecords.typeGeneral"),
  };

  return (
    <Tabs.Root value={typeFilter} className="mb-4">
      <Tabs.List className="flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
        {supportTypes.map((tab) => (
          <Tabs.Tab
            key={tab}
            value={tab}
            nativeButton={false}
            render={
              <Link
                to={
                  tab
                    ? "/projects/$projectId/support-records/$type"
                    : "/projects/$projectId/support-records"
                }
                params={tab ? { projectId, type: tab } : { projectId }}
              />
            }
            className="cursor-pointer rounded-sm px-4 py-1.5 m-0.5 text-sm font-medium text-fg-tertiary transition-colors hover:text-fg data-active:bg-bg data-active:text-fg data-active:shadow-card"
          >
            {tabLabels[tab]}
          </Tabs.Tab>
        ))}
      </Tabs.List>
    </Tabs.Root>
  );
}
