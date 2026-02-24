import { SignOutIcon } from "@/components/icons";
import {
  createFileRoute,
  Link,
  Navigate,
  Outlet,
} from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { AppFooter } from "@/components/app-footer";
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
    <div className="flex min-h-screen flex-col bg-bg">
      <header className="glass sticky top-0 z-50 border-b border-border-secondary">
        <div className="flex h-13 items-center justify-between px-5">
          <Link
            to="/"
            className="flex items-center gap-2.5 text-sm font-semibold text-fg hover:text-fg"
          >
            <span className="brand-icon inline-flex size-7 items-center justify-center rounded-lg text-xs font-bold text-white">
              O
            </span>
            {t("common.appName")}
          </Link>
          <div className="flex items-center gap-3">
            <span className="inline-flex size-7 items-center justify-center rounded-full bg-bg-tertiary text-[11px] font-semibold text-fg-secondary">
              {user?.email?.charAt(0).toUpperCase()}
            </span>
            <button
              type="button"
              onClick={() => logout()}
              className="flex cursor-pointer items-center gap-1.5 text-sm text-fg-tertiary hover:text-fg"
              title={t("auth.logout")}
            >
              <SignOutIcon size={16} />
            </button>
          </div>
        </div>
      </header>
      <div className="flex flex-1">
        <Outlet />
      </div>
      <AppFooter />
    </div>
  );
}
