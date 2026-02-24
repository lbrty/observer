export type ProjectRole = "owner" | "manager" | "consultant" | "viewer";

export interface ProjectPermission {
  id: string;
  project_id: string;
  user_id: string;
  role: ProjectRole;
  can_view_contact: boolean;
  can_view_personal: boolean;
  can_view_documents: boolean;
  created_at: string;
  updated_at: string;
}

export interface PermissionListOutput {
  permissions: ProjectPermission[];
}

export interface AssignPermissionInput {
  user_id: string;
  role: ProjectRole;
  can_view_contact: boolean;
  can_view_personal: boolean;
  can_view_documents: boolean;
}

export interface UpdatePermissionInput {
  role?: ProjectRole;
  can_view_contact?: boolean;
  can_view_personal?: boolean;
  can_view_documents?: boolean;
}
