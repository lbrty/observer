import { useMutation, useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import { useAuth } from "@/stores/auth";
import type { MFASetupData } from "@/types/auth";
import type { User } from "@/types/auth";

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
  const { setUser } = useAuth();
  return useMutation({
    mutationFn: ({ secret, totpCode }: { secret: string; totpCode: string }) =>
      api.post("auth/mfa/enable", { json: { secret, totp_code: totpCode } }),
    onSuccess: async () => {
      const user = await api.get("auth/me").json<User>();
      setUser(user);
    },
  });
}

export function useDisableMFA() {
  const { setUser } = useAuth();
  return useMutation({
    mutationFn: (totpCode: string) =>
      api.post("auth/mfa/disable", { json: { totp_code: totpCode } }),
    onSuccess: async () => {
      const user = await api.get("auth/me").json<User>();
      setUser(user);
    },
  });
}
