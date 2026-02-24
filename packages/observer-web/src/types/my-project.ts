import type { ProjectRole } from "./permission";

export interface MyProject {
  id: string;
  name: string;
  description?: string;
  status: string;
  role: string;
  can_view_contact: boolean;
  can_view_personal: boolean;
  can_view_documents: boolean;
  created_at: string;
  updated_at: string;
}

export interface MyProjectsOutput {
  projects: MyProject[];
}
