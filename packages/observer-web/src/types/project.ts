export type ProjectStatus = "active" | "archived" | "closed";

export interface Project {
  id: string;
  name: string;
  description?: string;
  owner_id: string;
  status: ProjectStatus;
  created_at: string;
  updated_at: string;
}

export interface ListProjectsOutput {
  projects: Project[];
  total: number;
  page: number;
  per_page: number;
}

export interface ListProjectsParams {
  owner_id?: string;
  status?: string;
  page?: number;
  per_page?: number;
}

export interface CreateProjectInput {
  name: string;
  description?: string;
}

export interface UpdateProjectInput {
  name?: string;
  description?: string;
  status?: ProjectStatus;
}
