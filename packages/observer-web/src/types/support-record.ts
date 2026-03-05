export type SupportType =
  | "humanitarian"
  | "legal"
  | "social"
  | "psychological"
  | "medical"
  | "general";

export type SupportSphere =
  | "housing_assistance"
  | "document_recovery"
  | "social_benefits"
  | "property_rights"
  | "employment_rights"
  | "family_law"
  | "healthcare_access"
  | "education_access"
  | "financial_aid"
  | "psychological_support"
  | "other";

export type ReferralStatus =
  | "pending"
  | "accepted"
  | "completed"
  | "declined"
  | "no_response";

export interface SupportRecord {
  id: string;
  person_id: string;
  project_id: string;
  consultant_id?: string;
  recorded_by?: string;
  office_id?: string;
  referred_to_office?: string;
  type: SupportType;
  sphere?: SupportSphere;
  referral_status?: ReferralStatus;
  provided_at?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateSupportRecordInput {
  person_id: string;
  type: SupportType;
  sphere?: SupportSphere;
  consultant_id?: string;
  office_id?: string;
  referred_to_office?: string;
  referral_status?: ReferralStatus;
  provided_at?: string;
  notes?: string;
}

export interface UpdateSupportRecordInput {
  type?: SupportType;
  sphere?: SupportSphere;
  consultant_id?: string;
  office_id?: string;
  referred_to_office?: string;
  referral_status?: ReferralStatus;
  provided_at?: string;
  notes?: string;
}

export interface ListSupportRecordsParams {
  person_id?: string;
  consultant_id?: string;
  office_id?: string;
  type?: SupportType;
  sphere?: SupportSphere;
  page?: number;
  per_page?: number;
}

export interface ListSupportRecordsOutput {
  support_records: SupportRecord[];
  total: number;
  page: number;
  per_page: number;
}
