import { GearSix, SignOut } from "@phosphor-icons/react";
import {
  createFileRoute,
  Link,
  Navigate,
  Outlet,
} from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_app")({
  component: AppLayout,
});

function AppLayout() {
  const { t } = useTranslation();
  const { isAuthenticated, isLoading, user, logout } = useAuth();

  if (isLoading) return null;

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  return (
    <div className="min-h-screen bg-bg">
      <header className="border-b border-border-secondary bg-bg-secondary">
        <div className="mx-auto flex h-14 max-w-5xl items-center justify-between px-4">
          <div className="flex items-center gap-4">
            <span className="text-sm font-semibold text-fg">
              {t("common.appName")}
            </span>
            {(user?.role === "admin" || user?.role === "staff") && (
              <Link
                to="/admin"
                className="flex items-center gap-1 text-sm text-fg-tertiary hover:text-fg"
              >
                <GearSix size={16} />
                {t("admin.title")}
              </Link>
            )}
          </div>
          <div className="flex items-center gap-3">
            <span className="text-sm text-fg-secondary">{user?.email}</span>
            <button
              type="button"
              onClick={() => logout()}
              className="flex cursor-pointer items-center gap-1 text-sm text-fg-tertiary hover:text-fg"
            >
              <SignOut size={16} />
              {t("common.logout")}
            </button>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-5xl px-4 py-6">
        <Outlet />
      </main>
    </div>
  );
}
