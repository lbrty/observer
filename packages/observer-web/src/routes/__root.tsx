import type { QueryClient } from "@tanstack/react-query";
import { Outlet, createRootRouteWithContext } from "@tanstack/react-router";

import { ToastContainer } from "@/components/toast";
import { AuthProvider } from "@/stores/auth";
import { ToastProvider } from "@/stores/toast";

interface RouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
});

function RootComponent() {
  return (
    <ToastProvider>
      <AuthProvider>
        <Outlet />
      </AuthProvider>
      <ToastContainer />
    </ToastProvider>
  );
}
