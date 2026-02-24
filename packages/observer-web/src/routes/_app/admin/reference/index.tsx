import { createFileRoute, Navigate } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/admin/reference/")({
  component: () => <Navigate to="/admin/reference/countries" />,
});
