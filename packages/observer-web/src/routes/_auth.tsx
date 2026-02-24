import { createFileRoute, Navigate, Outlet } from "@tanstack/react-router";

import { AppFooter } from "@/components/app-footer";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_auth")({
  component: AuthLayout,
});

function AuthLayout() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) return null;

  if (isAuthenticated) {
    return <Navigate to="/" />;
  }

  return (
    <div className="auth-backdrop flex min-h-screen flex-col bg-bg">
      <div className="relative flex flex-1 items-center justify-center px-4">
        <div className="w-full max-w-md">
          <div className="rounded-2xl border border-border-secondary bg-bg-secondary p-10 shadow-elevated">
            <Outlet />
          </div>
        </div>
      </div>
      <AppFooter />
    </div>
  );
}
