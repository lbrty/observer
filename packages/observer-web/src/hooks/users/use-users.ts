import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { UseQueryOptions } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  AdminUser,
  CreateUserInput,
  ListUsersOutput,
  ListUsersParams,
  UpdateUserInput,
} from "@/types/admin";

export function useUsers(
  params: ListUsersParams = {},
  options?: Partial<UseQueryOptions<ListUsersOutput>>,
) {
  return useQuery({
    queryKey: ["users", params],
    queryFn: ({ signal }) =>
      api
        .get("admin/users", { signal, searchParams: params as Record<string, string> })
        .json<ListUsersOutput>(),
    placeholderData: keepPreviousData,
    ...options,
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
    queryFn: ({ signal }) =>
      api
        .get("admin/users", {
          signal,
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

export function useDeactivateUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.patch(`admin/users/${id}/deactivate`).json<AdminUser>(),
    onSuccess: (_, id) => {
      qc.invalidateQueries({ queryKey: ["users"] });
      qc.invalidateQueries({ queryKey: ["users", id] });
    },
  });
}

export function useReactivateUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.patch(`admin/users/${id}/reactivate`).json<AdminUser>(),
    onSuccess: (_, id) => {
      qc.invalidateQueries({ queryKey: ["users"] });
      qc.invalidateQueries({ queryKey: ["users", id] });
    },
  });
}
