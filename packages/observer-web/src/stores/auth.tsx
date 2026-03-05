import { createContext, useContext, useEffect, useState } from "react";
import type { ReactNode } from "react";

import { api, HTTPError } from "@/lib/api";
import type { LoginInput, LoginOutput, RegisterInput, RegisterOutput, User } from "@/types/auth";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

interface AuthActions {
  login: (input: LoginInput) => Promise<LoginOutput>;
  register: (input: RegisterInput) => Promise<RegisterOutput>;
  logout: () => Promise<void>;
  setUser: (user: User) => void;
}

type AuthContextValue = AuthState & AuthActions;

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = user !== null;

  useEffect(() => {
    api
      .get("auth/me")
      .json<User>()
      .then(setUser)
      .catch(() => {})
      .finally(() => setIsLoading(false));
  }, []);

  async function login(input: LoginInput): Promise<LoginOutput> {
    const data = await api.post("auth/login", { json: input }).json<LoginOutput>();

    if (!data.requires_mfa && data.user) {
      setUser(data.user);
    }

    return data;
  }

  async function register(input: RegisterInput): Promise<RegisterOutput> {
    return api.post("auth/register", { json: input }).json<RegisterOutput>();
  }

  async function logout() {
    try {
      await api.post("auth/logout");
    } catch (err) {
      if (!(err instanceof HTTPError)) throw err;
    }
    setUser(null);
  }

  return (
    <AuthContext
      value={{
        user,
        isAuthenticated,
        isLoading,
        login,
        register,
        logout,
        setUser,
      }}
    >
      {children}
    </AuthContext>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
}
