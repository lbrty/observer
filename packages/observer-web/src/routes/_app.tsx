import { SignOut } from "@phosphor-icons/react";
import { createFileRoute, Navigate, Outlet } from "@tanstack/react-router";
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
    <div className="min-h-screen bg-zinc-50">
      <header className="border-b border-zinc-200 bg-white">
        <div className="mx-auto flex h-14 max-w-5xl items-center justify-between px-4">
          <span className="text-sm font-semibold">{t("common.appName")}</span>
          <div className="flex items-center gap-3">
            <span className="text-sm text-zinc-600">{user?.email}</span>
            <button
              type="button"
              onClick={() => logout()}
              className="flex cursor-pointer items-center gap-1 text-sm text-zinc-500 hover:text-zinc-900"
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
