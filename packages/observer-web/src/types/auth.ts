export type Role = "admin" | "staff" | "consultant" | "guest";

export interface User {
  id: string;
  email: string;
  phone: string;
  role: Role;
  is_verified: boolean;
  created_at: string;
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
  phone: string;
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

export interface ApiError {
  error: string;
}
