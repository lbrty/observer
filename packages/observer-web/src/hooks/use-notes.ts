import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { CreateNoteInput, ListNotesOutput, Note } from "@/types/note";

export function useNotes(projectId: string, personId: string) {
  return useQuery({
    queryKey: ["notes", projectId, personId],
    queryFn: () =>
      api.get(`projects/${projectId}/people/${personId}/notes`).json<ListNotesOutput>(),
    enabled: !!projectId && !!personId,
  });
}

export function useCreateNote(projectId: string, personId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateNoteInput) =>
      api.post(`projects/${projectId}/people/${personId}/notes`, { json: data }).json<Note>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["notes", projectId, personId] }),
  });
}

export function useDeleteNote(projectId: string, personId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`projects/${projectId}/people/${personId}/notes/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["notes", projectId, personId] }),
  });
}
