export type CaseStatus = "new" | "active" | "closed" | "archived";
export type Sex = "male" | "female" | "other" | "unknown";
export type AgeGroup =
  | "infant"
  | "toddler"
  | "pre_school"
  | "middle_childhood"
  | "young_teen"
  | "teenager"
  | "young_adult"
  | "early_adult"
  | "middle_aged_adult"
  | "old_adult"
  | "unknown";

export interface Person {
  id: string;
  project_id: string;
  consultant_id?: string;
  office_id?: string;
  current_place_id?: string;
  origin_place_id?: string;
  external_id?: string;
  first_name: string;
  last_name?: string;
  patronymic?: string;
  email?: string;
  birth_date?: string;
  sex: Sex;
  age_group?: AgeGroup;
  primary_phone?: string;
  phone_numbers: string[];
  case_status: CaseStatus;
  consent_given: boolean;
  consent_date?: string;
  registered_at?: string;
  created_at: string;
  updated_at: string;
}

export interface ListPeopleOutput {
  people: Person[];
  total: number;
  page: number;
  per_page: number;
}

export interface ListPeopleParams {
  search?: string;
  case_status?: string;
  office_id?: string;
  consultant_id?: string;
  page?: number;
  per_page?: number;
}

export interface CreatePersonInput {
  consultant_id?: string;
  office_id?: string;
  current_place_id?: string;
  origin_place_id?: string;
  external_id?: string;
  first_name: string;
  last_name?: string;
  patronymic?: string;
  email?: string;
  birth_date?: string;
  sex?: Sex;
  age_group?: AgeGroup;
  primary_phone?: string;
  phone_numbers?: string[];
  case_status?: CaseStatus;
  consent_given?: boolean;
  consent_date?: string;
  registered_at?: string;
}

export interface UpdatePersonInput {
  consultant_id?: string;
  office_id?: string;
  current_place_id?: string;
  origin_place_id?: string;
  external_id?: string;
  first_name?: string;
  last_name?: string;
  patronymic?: string;
  email?: string;
  birth_date?: string;
  sex?: Sex;
  age_group?: AgeGroup;
  primary_phone?: string;
  phone_numbers?: string[];
  case_status?: CaseStatus;
  consent_given?: boolean;
  consent_date?: string;
  registered_at?: string;
}
