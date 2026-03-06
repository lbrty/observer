import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { EmptyState } from "@/components/empty-state";
import { FilesIcon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";

export const Route = createFileRoute("/_app/projects/$projectId/documents/")({
  component: DocumentsPage,
});

function DocumentsPage() {
  const { t } = useTranslation();

  return (
    <div>
      <PageHeader title={t("project.documents.title")} />
      <div className="rounded-xl border border-border-secondary bg-bg-secondary">
        <EmptyState
          icon={FilesIcon}
          title={t("project.documents.emptyTitle")}
          description={t("project.documents.emptyDescription")}
        />
      </div>
    </div>
  );
}
