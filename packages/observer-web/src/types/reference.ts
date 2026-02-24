export interface Country {
  id: string;
  name: string;
  code: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCountryInput {
  name: string;
  code: string;
}

export interface UpdateCountryInput {
  name?: string;
  code?: string;
}

export interface State {
  id: string;
  country_id: string;
  name: string;
  code?: string;
  conflict_zone?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateStateInput {
  name: string;
  code?: string;
  conflict_zone?: string;
}

export interface UpdateStateInput {
  name?: string;
  code?: string;
  conflict_zone?: string;
}

export interface Place {
  id: string;
  state_id: string;
  name: string;
  lat?: number;
  lon?: number;
  created_at: string;
  updated_at: string;
}

export interface CreatePlaceInput {
  name: string;
  lat?: number;
  lon?: number;
}

export interface UpdatePlaceInput {
  name?: string;
  lat?: number;
  lon?: number;
}

export interface Office {
  id: string;
  name: string;
  place_id?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateOfficeInput {
  name: string;
  place_id?: string;
}

export interface UpdateOfficeInput {
  name?: string;
  place_id?: string;
}

export interface Category {
  id: string;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCategoryInput {
  name: string;
  description?: string;
}

export interface UpdateCategoryInput {
  name?: string;
  description?: string;
}
