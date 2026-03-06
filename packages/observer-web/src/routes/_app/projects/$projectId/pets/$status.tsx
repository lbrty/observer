import { createFileRoute, useNavigate } from "@tanstack/react-router";

import { PetsContent, type PetStatus } from "./-pets-page";

const validStatuses = new Set<string>([
  "registered",
  "adopted",
  "owner_found",
  "needs_shelter",
  "unknown",
]);

export const Route = createFileRoute("/_app/projects/$projectId/pets/$status")({
  component: PetsByStatus,
  validateSearch: (search: Record<string, unknown>): { page?: number } => ({
    page: Number(search.page) || undefined,
  }),
});

function PetsByStatus() {
  const { projectId, status } = Route.useParams();
  const navigate = useNavigate();
  const { page = 1 } = Route.useSearch();

  const statusFilter: PetStatus = validStatuses.has(status) ? (status as PetStatus) : "";

  function setPage(value: number) {
    navigate({ from: Route.fullPath, search: { page: value > 1 ? value : undefined }, replace: true });
  }

  return (
    <PetsContent
      projectId={projectId}
      statusFilter={statusFilter}
      page={page}
      onPageChange={setPage}
    />
  );
}
