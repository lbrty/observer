export interface AuditEntry {
  id: string;
  project_id: string | null;
  user_id: string | null;
  action: string;
  entity_type: string;
  entity_id: string | null;
  summary: string;
  details: Record<string, unknown> | null;
  ip: string | null;
  user_agent: string | null;
  created_at: string;
  user_first_name: string | null;
  user_last_name: string | null;
  user_email: string | null;
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
