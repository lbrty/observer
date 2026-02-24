import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { MyProjectsOutput } from "@/types/my-project";

export function useMyProjects() {
  return useQuery({
    queryKey: ["my-projects"],
    queryFn: () => api.get("my/projects").json<MyProjectsOutput>(),
  });
}
