import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { PageHeader } from "@/components/page-header";

export const Route = createFileRoute("/_app/projects/$projectId/documents/")({
  component: DocumentsPage,
});

function DocumentsPage() {
  const { t } = useTranslation();

  return (
    <div>
      <PageHeader title={t("project.documents.title")} />
      <div className="rounded-xl border border-border-secondary bg-bg-secondary p-8 text-center text-sm text-fg-tertiary">
        {t("project.documents.empty")}
      </div>
    </div>
  );
}
