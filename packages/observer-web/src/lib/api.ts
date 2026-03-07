import ky, { HTTPError } from "ky";

const BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:9000";

let refreshPromise: Promise<void> | null = null;

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
          return response;
        }
      },
    ],
  },
});

export { HTTPError };
