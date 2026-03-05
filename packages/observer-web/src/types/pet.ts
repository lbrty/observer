export type PetStatus =
  | "registered"
  | "adopted"
  | "owner_found"
  | "needs_shelter"
  | "unknown";

export interface Pet {
  id: string;
  project_id: string;
  owner_id?: string;
  name: string;
  status: PetStatus;
  registration_id?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePetInput {
  owner_id?: string;
  name: string;
  status?: PetStatus;
  registration_id?: string;
  notes?: string;
}

export interface UpdatePetInput {
  owner_id?: string;
  name?: string;
  status?: PetStatus;
  registration_id?: string;
  notes?: string;
}

export interface ListPetsParams {
  page?: number;
  per_page?: number;
}

export interface ListPetsOutput {
  pets: Pet[];
  total: number;
  page: number;
  per_page: number;
}
