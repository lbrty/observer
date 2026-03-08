import { createFileRoute, useNavigate } from "@tanstack/react-router";

import { SupportRecordsContent } from "./-support-records-page";

export const Route = createFileRoute("/_app/projects/$projectId/support-records/")({
  component: SupportRecordsPage,
  validateSearch: (search: Record<string, unknown>): { page?: number } => ({
    page: Number(search.page) || undefined,
  }),
});

function SupportRecordsPage() {
  const { projectId } = Route.useParams();
  const navigate = useNavigate();
  const { page = 1 } = Route.useSearch();

  function setPage(value: number) {
    navigate({
      from: Route.fullPath,
      search: { page: value > 1 ? value : undefined },
      replace: true,
    });
  }

  return (
    <SupportRecordsContent projectId={projectId} typeFilter="" page={page} onPageChange={setPage} />
  );
}
