import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/admin/reference/countries/$countryId")({
  component: () => <Outlet />,
});
