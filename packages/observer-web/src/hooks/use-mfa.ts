import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { MFASetupData } from "@/types/auth";

export function useMFASetup(enabled: boolean) {
  return useQuery({
    queryKey: ["mfa-setup"],
    queryFn: () => api.get("auth/mfa/setup").json<MFASetupData>(),
    enabled,
    staleTime: 0,
    gcTime: 0,
  });
}

export function useEnableMFA() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ secret, totpCode }: { secret: string; totpCode: string }) =>
      api.post("auth/mfa/enable", { json: { secret, totp_code: totpCode } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["me"] }),
  });
}

export function useDisableMFA() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (totpCode: string) =>
      api.post("auth/mfa/disable", { json: { totp_code: totpCode } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["me"] }),
  });
}
