import { SignOutIcon, UserCircleIcon } from "@/components/icons";
import { Menu } from "@base-ui/react/menu";
import { createFileRoute, Link, Navigate, Outlet } from "@tanstack/react-router";
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
          <AvatarMenu email={user?.email ?? ""} onLogout={logout} />
        </div>
      </header>
      <div className="flex flex-1">
        <Outlet />
      </div>
    </div>
  );
}

function AvatarMenu({ email, onLogout }: { email: string; onLogout: () => void }) {
  const { t } = useTranslation();

  return (
    <Menu.Root>
      <Menu.Trigger className="inline-flex size-7 cursor-pointer items-center justify-center rounded-full bg-bg-tertiary text-[11px] font-semibold text-fg-secondary transition-shadow hover:ring-2 hover:ring-accent/30">
        {email.charAt(0).toUpperCase()}
      </Menu.Trigger>
      <Menu.Portal>
        <Menu.Positioner sideOffset={6} align="end" className="z-[100]">
          <Menu.Popup className="w-44 origin-(--transform-origin) rounded-xl border border-border-secondary bg-bg-secondary py-1 shadow-elevated transition-[transform,scale,opacity] data-ending-style:scale-95 data-ending-style:opacity-0 data-starting-style:scale-95 data-starting-style:opacity-0">
            <Menu.Item
              render={<Link to="/profile" />}
              className="flex cursor-pointer items-center gap-2 px-3 py-1.5 text-sm text-fg outline-none select-none data-highlighted:bg-bg-tertiary"
            >
              <span className="inline-flex w-4 items-center justify-center text-fg-tertiary">
                <UserCircleIcon size={14} />
              </span>
              {t("profile.title")}
            </Menu.Item>

            <Menu.Separator className="my-1 h-px bg-border-secondary" />

            <Menu.Item
              onClick={onLogout}
              className="flex cursor-pointer items-center gap-2 px-3 py-1.5 text-sm text-fg outline-none select-none data-highlighted:bg-bg-tertiary"
            >
              <span className="inline-flex w-4 items-center justify-center text-fg-tertiary">
                <SignOutIcon size={14} />
              </span>
              {t("common.logout")}
            </Menu.Item>
          </Menu.Popup>
        </Menu.Positioner>
      </Menu.Portal>
    </Menu.Root>
  );
}
