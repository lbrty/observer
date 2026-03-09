import { useMyProjects } from "./use-my-projects";

export function useProjectRole(projectId: string) {
  const { data } = useMyProjects();
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
