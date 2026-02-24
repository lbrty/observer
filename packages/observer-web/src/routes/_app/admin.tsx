import {
  Buildings,
  FolderSimple,
  Globe,
  MapPin,
  MapPinArea,
  Tag,
  Users,
} from "@phosphor-icons/react";
import { createFileRoute, Navigate, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { SidebarLink } from "@/components/sidebar-link";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_app/admin")({
  component: AdminLayout,
});

function AdminLayout() {
  const { t } = useTranslation();
  const { user } = useAuth();

  const isAdmin = user?.role === "admin";
  const isStaff = user?.role === "staff";

  if (!isAdmin && !isStaff) {
    return <Navigate to="/" />;
  }

  return (
    <div className="flex gap-6">
      <aside className="w-60 shrink-0">
        <nav className="sticky top-20 space-y-1">
          <SidebarLink
            to="/admin/users"
            label={t("admin.nav.users")}
            icon={Users}
          />

          {isAdmin && (
            <SidebarLink
              to="/admin/projects"
              label={t("admin.nav.projects")}
              icon={FolderSimple}
            />
          )}

          <div className="pt-4 pb-1 pl-3 text-xs font-medium uppercase tracking-wider text-fg-tertiary">
            {t("admin.nav.reference")}
          </div>
          <SidebarLink
            to="/admin/reference/countries"
            label={t("admin.nav.countries")}
            icon={Globe}
          />
          <SidebarLink
            to="/admin/reference/states"
            label={t("admin.nav.states")}
            icon={MapPin}
          />
          <SidebarLink
            to="/admin/reference/places"
            label={t("admin.nav.places")}
            icon={MapPinArea}
          />
          <SidebarLink
            to="/admin/reference/offices"
            label={t("admin.nav.offices")}
            icon={Buildings}
          />
          <SidebarLink
            to="/admin/reference/categories"
            label={t("admin.nav.categories")}
            icon={Tag}
          />
        </nav>
      </aside>
      <div className="min-w-0 flex-1">
        <Outlet />
      </div>
    </div>
  );
}
