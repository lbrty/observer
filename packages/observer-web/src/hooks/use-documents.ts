import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  Document,
  ListDocumentsOutput,
  UpdateDocumentInput,
} from "@/types/document";

export function useDocuments(projectId: string, personId: string) {
  return useQuery({
    queryKey: ["documents", projectId, personId],
    queryFn: () =>
      api.get(`projects/${projectId}/people/${personId}/documents`).json<ListDocumentsOutput>(),
    enabled: !!projectId && !!personId,
  });
}

export function useDocument(projectId: string, id: string) {
  return useQuery({
    queryKey: ["documents", projectId, id],
    queryFn: () => api.get(`projects/${projectId}/documents/${id}`).json<Document>(),
    enabled: !!projectId && !!id,
  });
}

export function useUploadDocument(projectId: string, personId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (file: File) => {
      const formData = new FormData();
      formData.append("file", file);
      return api
        .post(`projects/${projectId}/people/${personId}/documents`, { body: formData })
        .json<Document>();
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["documents", projectId] }),
  });
}

export function useUpdateDocument(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateDocumentInput }) =>
      api.patch(`projects/${projectId}/documents/${id}`, { json: data }).json<Document>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["documents", projectId] }),
  });
}

export function useDeleteDocument(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`projects/${projectId}/documents/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["documents", projectId] }),
  });
}

export function documentDownloadUrl(projectId: string, documentId: string): string {
  const base = import.meta.env.VITE_API_URL ?? "http://localhost:9000";
  return `${base}/projects/${projectId}/documents/${documentId}/download`;
}

export function documentStreamUrl(projectId: string, documentId: string): string {
  const base = import.meta.env.VITE_API_URL ?? "http://localhost:9000";
  return `${base}/projects/${projectId}/documents/${documentId}/stream`;
}

export function documentThumbnailUrl(projectId: string, documentId: string): string {
  const base = import.meta.env.VITE_API_URL ?? "http://localhost:9000";
  return `${base}/projects/${projectId}/documents/${documentId}/thumbnail`;
}

export function isImageMime(mimeType: string): boolean {
  return mimeType.startsWith("image/");
}

export function isPdfMime(mimeType: string): boolean {
  return mimeType === "application/pdf";
}
