import { createLazyFileRoute, Navigate } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/_app/admin/")({
  component: () => <Navigate to="/admin/users" />,
});
