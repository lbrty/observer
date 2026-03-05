import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  AdminUser,
  CreateUserInput,
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
    placeholderData: keepPreviousData,
  });
}

export function useUser(id: string) {
  return useQuery({
    queryKey: ["users", id],
    queryFn: () => api.get(`admin/users/${id}`).json<AdminUser>(),
    enabled: !!id,
  });
}

export function useSearchUsers(search: string) {
  return useQuery({
    queryKey: ["users", "search", search],
    queryFn: () =>
      api
        .get("admin/users", {
          searchParams: { search, per_page: "10" },
        })
        .json<ListUsersOutput>(),
    enabled: search.length >= 2,
    placeholderData: keepPreviousData,
  });
}

export function useCreateUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateUserInput) =>
      api.post("admin/users", { json: data }).json<AdminUser>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["users"] }),
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
