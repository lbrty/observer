import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { SearchOutput } from "@/types/search";

export function useSearch(query: string, limit = 5) {
  return useQuery({
    queryKey: ["search", query, limit],
    queryFn: ({ signal }) =>
      api
        .get(`search?q=${encodeURIComponent(query)}&limit=${limit}`, { signal })
        .json<SearchOutput>(),
    enabled: query.length >= 2,
    staleTime: 30_000,
  });
}
