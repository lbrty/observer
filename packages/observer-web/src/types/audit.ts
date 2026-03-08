export interface AuditEntry {
  id: string;
  project_id: string | null;
  user_id: string | null;
  action: string;
  entity_type: string;
  entity_id: string | null;
  summary: string;
  ip: string;
  user_agent: string;
  created_at: string;
  user_first_name: string;
  user_last_name: string;
  user_email: string;
}

export interface AuditListParams {
  project_id?: string;
  user_id?: string;
  action?: string;
  entity_type?: string;
  date_from?: string;
  date_to?: string;
  page?: number;
  per_page?: number;
}

export interface AuditListOutput {
  entries: AuditEntry[];
  total: number;
  page: number;
  per_page: number;
}
