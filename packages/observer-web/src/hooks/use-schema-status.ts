import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import { useAuth } from "@/stores/auth";

interface SchemaStatus {
  current_version: number;
  latest_version: number;
  pending: number;
  dirty: boolean;
}

export function useSchemaStatus() {
  const { user } = useAuth();
  return useQuery({
    queryKey: ["schema-status"],
    queryFn: () => api.get("admin/schema/status").json<SchemaStatus>(),
    enabled: user?.role === "admin",
    staleTime: 5 * 60 * 1000,
  });
}
