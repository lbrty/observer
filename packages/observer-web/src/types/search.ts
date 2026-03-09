export interface PersonResult {
  id: string;
  first_name: string;
  last_name: string;
}

export interface PetResult {
  id: string;
  name: string;
}

export interface ProjectResult {
  id: string;
  name: string;
}

export interface ProjectGroup {
  project_id: string;
  project_name: string;
  people: PersonResult[];
  pets: PetResult[];
  projects: ProjectResult[];
}

export interface SearchOutput {
  query: string;
  results: ProjectGroup[];
  total: number;
}
