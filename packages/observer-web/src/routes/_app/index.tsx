import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_app/")({
  component: DashboardPage,
});

function DashboardPage() {
  const { t } = useTranslation();
  const { user } = useAuth();

  return (
    <div>
      <h1 className="text-xl font-semibold">
        {t("dashboard.greeting", { email: user?.email })}
      </h1>
      <p className="mt-1 text-sm text-zinc-500">
        {t("dashboard.role", { role: user?.role })}
      </p>
    </div>
  );
}
