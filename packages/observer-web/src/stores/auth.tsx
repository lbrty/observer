import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import type { ReactNode } from "react";

import { api, ApiRequestError } from "@/lib/api";
import type {
  LoginInput,
  LoginOutput,
  RegisterInput,
  RegisterOutput,
  User,
} from "@/types/auth";

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
    api<User>("/auth/me")
      .then(setUser)
      .catch(() => {})
      .finally(() => setIsLoading(false));
  }, []);

  const login = useCallback(async (input: LoginInput): Promise<LoginOutput> => {
    const data = await api<LoginOutput>("/auth/login", {
      method: "POST",
      body: JSON.stringify(input),
    });

    if (!data.requires_mfa && data.user) {
      setUser(data.user);
    }

    return data;
  }, []);

  const register = useCallback(
    async (input: RegisterInput): Promise<RegisterOutput> => {
      return api<RegisterOutput>("/auth/register", {
        method: "POST",
        body: JSON.stringify(input),
      });
    },
    [],
  );

  const logout = useCallback(async () => {
    try {
      await api("/auth/logout", { method: "POST" });
    } catch (err) {
      if (!(err instanceof ApiRequestError)) throw err;
    }
    setUser(null);
  }, []);

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
