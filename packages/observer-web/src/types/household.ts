export type Relationship =
  | "head"
  | "spouse"
  | "child"
  | "parent"
  | "sibling"
  | "grandchild"
  | "grandparent"
  | "other_relative"
  | "non_relative";

export interface HouseholdMember {
  person_id: string;
  relationship: Relationship;
}

export interface Household {
  id: string;
  project_id: string;
  reference_number?: string;
  head_person_id?: string;
  member_count: number;
  members?: HouseholdMember[];
  created_at: string;
  updated_at: string;
}

export interface CreateHouseholdInput {
  reference_number?: string;
  head_person_id?: string;
}

export interface UpdateHouseholdInput {
  reference_number?: string;
  head_person_id?: string;
}

export interface AddMemberInput {
  person_id: string;
  relationship: Relationship;
}

export interface ListHouseholdsParams {
  page?: number;
  per_page?: number;
}

export interface ListHouseholdsOutput {
  households: Household[];
  total: number;
  page: number;
  per_page: number;
}
