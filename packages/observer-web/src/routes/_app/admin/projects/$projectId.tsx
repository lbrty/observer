import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/admin/projects/$projectId")({
  component: () => <Outlet />,
});
