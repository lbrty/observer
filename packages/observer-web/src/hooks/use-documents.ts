import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreateDocumentInput,
  Document,
  ListDocumentsOutput,
} from "@/types/document";

export function useDocuments(projectId: string, personId: string) {
  return useQuery({
    queryKey: ["documents", projectId, personId],
    queryFn: () =>
      api
        .get(`projects/${projectId}/people/${personId}/documents`)
        .json<ListDocumentsOutput>(),
    enabled: !!projectId && !!personId,
  });
}

export function useDocument(projectId: string, id: string) {
  return useQuery({
    queryKey: ["documents", projectId, id],
    queryFn: () =>
      api.get(`projects/${projectId}/documents/${id}`).json<Document>(),
    enabled: !!projectId && !!id,
  });
}

export function useCreateDocument(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateDocumentInput) =>
      api
        .post(`projects/${projectId}/documents`, { json: data })
        .json<Document>(),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["documents", projectId] }),
  });
}

export function useDeleteDocument(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) =>
      api.delete(`projects/${projectId}/documents/${id}`),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["documents", projectId] }),
  });
}
