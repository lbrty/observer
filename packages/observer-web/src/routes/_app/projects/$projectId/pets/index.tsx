import { createFileRoute, useNavigate } from "@tanstack/react-router";

import { PetsContent } from "./-pets-page";

export const Route = createFileRoute("/_app/projects/$projectId/pets/")({
  component: PetsPage,
  validateSearch: (search: Record<string, unknown>): { page?: number } => ({
    page: Number(search.page) || undefined,
  }),
});

function PetsPage() {
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

  return <PetsContent projectId={projectId} statusFilter="" page={page} onPageChange={setPage} />;
}
