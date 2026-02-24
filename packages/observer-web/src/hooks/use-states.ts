import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreateStateInput,
  State,
  UpdateStateInput,
} from "@/types/reference";

import { makeReferenceHooks } from "./make-reference-hooks";

const { useUpdate: useUpdateState, useDelete: useDeleteState } =
  makeReferenceHooks<State, CreateStateInput, UpdateStateInput>("states");

export { useUpdateState, useDeleteState };

export function useStates(countryId?: string) {
  return useQuery({
    queryKey: ["states", countryId ?? "all"],
    queryFn: () =>
      api
        .get("admin/states", {
          searchParams: countryId ? { country_id: countryId } : {},
        })
        .json<{ states: State[] }>(),
  });
}

export function useCreateState() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      countryId,
      data,
    }: {
      countryId: string;
      data: CreateStateInput;
    }) =>
      api
        .post("admin/states", {
          json: data,
          searchParams: { country_id: countryId },
        })
        .json<State>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["states"] }),
  });
}
