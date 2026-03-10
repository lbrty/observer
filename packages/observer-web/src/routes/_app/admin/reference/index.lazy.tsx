import { createLazyFileRoute, Navigate } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/_app/admin/reference/")({
  component: () => <Navigate to="/admin/reference/countries" />,
});
