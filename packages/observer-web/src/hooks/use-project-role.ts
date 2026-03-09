import { useAuth } from "@/stores/auth";

import { useMyProjects } from "./use-my-projects";

export function useProjectRole(projectId: string) {
  const { user } = useAuth();
  const { data } = useMyProjects();

  // Platform guests are always read-only regardless of project role assignment.
  if (user?.role === "guest") {
    return { role: "viewer" as const, canWrite: false, canDelete: false, canExport: false };
  }

  const project = data?.projects.find((p) => p.id === projectId);
  const role = project?.role;

  // consultant+ can create and update; viewer is read-only
  const canWrite = role === "owner" || role === "manager" || role === "consultant";
  // manager+ can delete
  const canDelete = role === "owner" || role === "manager";
  // owners and managers always have export; consultants need the explicit flag
  const canExport =
    role === "owner" ||
    role === "manager" ||
    (role === "consultant" && (project?.can_export ?? false));

  return { role, canWrite, canDelete, canExport };
}
