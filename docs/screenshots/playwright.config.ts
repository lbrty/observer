import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: ".",
  testMatch: "capture.ts",
  timeout: 120_000,
  use: {
    baseURL: "http://localhost:5173",
    viewport: { width: 1440, height: 900 },
    colorScheme: "light",
  },
  projects: [
    {
      name: "screenshots",
      use: { browserName: "chromium" },
    },
  ],
});
