import type { Role } from "./auth";

export interface AdminUser {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  office_id?: string;
  role: Role;
  is_verified: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ListUsersOutput {
  users: AdminUser[];
  total: number;
  page: number;
  per_page: number;
}

export interface ListUsersParams {
  page?: number;
  per_page?: number;
  search?: string;
  role?: string;
  is_active?: boolean;
}

export interface CreateUserInput {
  first_name: string;
  last_name?: string;
  email: string;
  phone?: string;
  password: string;
  role: string;
  office_id?: string | null;
  is_active: boolean;
  is_verified: boolean;
}

export interface UpdateUserInput {
  first_name?: string;
  last_name?: string;
  email?: string;
  phone?: string;
  office_id?: string | null;
  role?: string;
  is_active?: boolean;
  is_verified?: boolean;
}
