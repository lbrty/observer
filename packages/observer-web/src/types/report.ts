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
  people_by_sphere: ReportGroup;
  by_office: ReportGroup;
  by_age_group: ReportGroup;
  consultations_by_age_group: ReportGroup;
  by_tag: ReportGroup;
  family_units: ReportGroup;
  by_case_status: ReportGroup;
  status_flow: StatusFlow[];
}

export interface MonthlyStatusCount {
  month: string;
  status: string;
  count: number;
}

export interface PetReport {
  by_status: ReportGroup;
  by_ownership: ReportGroup;
  by_month: ReportGroup;
  by_status_by_month: MonthlyStatusCount[];
}

export interface PetReportParams {
  date_from?: string;
  date_to?: string;
  status?: string;
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
  support_type?: string;
}

export interface CustomReportParams {
  metric: "events" | "people" | "units" | "pets";
  group_by: string[];
  date_from?: string;
  date_to?: string;
  support_type?: string;
}

export interface CustomRow {
  dimensions: Record<string, string>;
  count: number;
}

export interface CustomReportOutput {
  metric: string;
  group_by: string[];
  rows: CustomRow[];
  total: number;
}
