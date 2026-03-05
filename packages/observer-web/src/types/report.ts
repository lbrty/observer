export interface CountResult {
  label: string;
  count: number;
}

export interface ReportGroup {
  group: string;
  rows: CountResult[];
  total: number;
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
}

export interface ReportParams {
  date_from?: string;
  date_to?: string;
}
