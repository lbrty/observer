export type Role = "admin" | "staff" | "consultant" | "guest";

export interface User {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  role: Role;
  is_verified: boolean;
  office_id?: string;
  created_at: string;
}

export interface UpdateProfileInput {
  first_name?: string;
  last_name?: string;
  phone?: string;
}

export interface ChangePasswordInput {
  current_password: string;
  new_password: string;
}

export interface ResetPasswordInput {
  new_password: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_at: string;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface LoginOutput {
  requires_mfa: boolean;
  mfa_token: string | null;
  tokens: TokenPair | null;
  user: User | null;
}

export interface RegisterInput {
  email: string;
  password: string;
  role: Role;
}

export interface RegisterOutput {
  user_id: string;
  message: string;
}

export interface RefreshInput {
  refresh_token: string;
}

export interface RefreshOutput {
  access_token: string;
  refresh_token: string;
  expires_at: string;
}

export interface LogoutInput {
  refresh_token: string;
}
