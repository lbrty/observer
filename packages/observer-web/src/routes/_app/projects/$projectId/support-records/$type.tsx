import { createFileRoute, useNavigate } from "@tanstack/react-router";

import { SupportRecordsContent, type SupportType } from "./-support-records-page";

const validTypes = new Set<string>([
  "humanitarian",
  "legal",
  "social",
  "psychological",
  "medical",
  "general",
]);

export const Route = createFileRoute("/_app/projects/$projectId/support-records/$type")({
  component: SupportRecordsByType,
  validateSearch: (search: Record<string, unknown>): { page?: number } => ({
    page: Number(search.page) || undefined,
  }),
});

function SupportRecordsByType() {
  const { projectId, type } = Route.useParams();
  const navigate = useNavigate();
  const { page = 1 } = Route.useSearch();

  const typeFilter: SupportType = validTypes.has(type) ? (type as SupportType) : "";

  function setPage(value: number) {
    navigate({ from: Route.fullPath, search: { page: value > 1 ? value : undefined }, replace: true });
  }

  return (
    <SupportRecordsContent
      projectId={projectId}
      typeFilter={typeFilter}
      page={page}
      onPageChange={setPage}
    />
  );
}
