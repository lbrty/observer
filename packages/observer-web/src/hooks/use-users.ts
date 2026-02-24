import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  AdminUser,
  ListUsersOutput,
  ListUsersParams,
  UpdateUserInput,
} from "@/types/admin";

export function useUsers(params: ListUsersParams = {}) {
  return useQuery({
    queryKey: ["users", params],
    queryFn: () =>
      api
        .get("admin/users", { searchParams: params as Record<string, string> })
        .json<ListUsersOutput>(),
  });
}

export function useUser(id: string) {
  return useQuery({
    queryKey: ["users", id],
    queryFn: () => api.get(`admin/users/${id}`).json<AdminUser>(),
    enabled: !!id,
  });
}

export function useUpdateUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserInput }) =>
      api.patch(`admin/users/${id}`, { json: data }).json<AdminUser>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["users"] }),
  });
}
