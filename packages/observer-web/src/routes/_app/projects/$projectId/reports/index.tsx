import { createFileRoute, Navigate } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/projects/$projectId/reports/")({
  component: ReportsIndex,
});

function ReportsIndex() {
  const { projectId } = Route.useParams();
  return <Navigate to="/projects/$projectId/reports/people" params={{ projectId }} replace />;
}
