export interface AuditEntry {
  id: string;
  project_id: string | null;
  user_id: string;
  action: string;
  entity_type: string;
  entity_id: string | null;
  summary: string;
  ip: string;
  user_agent: string;
  created_at: string;
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
