import ky, { HTTPError } from "ky";

export const apiBase = import.meta.env.VITE_API_URL ?? "http://localhost:9000";

const BASE_URL = apiBase;

let refreshPromise: Promise<void> | null = null;

// Registered by AuthProvider to clear auth state when session cannot be recovered.
let onUnauthorized: (() => void) | null = null;
export function setOnUnauthorized(fn: () => void) {
  onUnauthorized = fn;
}

async function refreshTokens(): Promise<void> {
  const res = await fetch(`${BASE_URL}/auth/refresh`, {
    method: "POST",
    credentials: "include",
  });

  if (!res.ok) {
    throw new Error("refresh failed");
  }
}

export const api = ky.create({
  prefixUrl: BASE_URL,
  credentials: "include",
  hooks: {
    beforeRequest: [
      (request) => {
        const csrf = document.cookie.match(/csrf_token=([^;]+)/)?.[1];
        if (csrf) request.headers.set("X-CSRF-Token", csrf);
      },
    ],
    afterResponse: [
      async (request, _options, response) => {
        if (response.status !== 401) return response;
        if (request.url.includes("/auth/login")) return response;

        try {
          if (!refreshPromise) {
            refreshPromise = refreshTokens();
          }
          await refreshPromise;
          refreshPromise = null;
          return ky(request, { credentials: "include" });
        } catch {
          refreshPromise = null;
          onUnauthorized?.();
          return response;
        }
      },
    ],
  },
});

export { HTTPError };
