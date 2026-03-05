export interface CountResult {
  label: string;
  count: number;
}

export interface ReportGroup {
  group: string;
  rows: CountResult[];
  total: number;
}

export interface StatusFlow {
  from_status: string;
  to_status: string;
  count: number;
  avg_days: number;
}

export interface FullReport {
  consultations: ReportGroup;
  by_sex: ReportGroup;
  by_idp_status: ReportGroup;
  by_category: ReportGroup;
  by_region: ReportGroup;
  by_sphere: ReportGroup;
  by_office: ReportGroup;
  by_age_group: ReportGroup;
  by_tag: ReportGroup;
  family_units: ReportGroup;
  status_flow: StatusFlow[];
}

export interface ReportParams {
  date_from?: string;
  date_to?: string;
  office_id?: string;
  category_id?: string;
  consultant_id?: string;
  case_status?: string;
  sex?: string;
  age_group?: string;
}
