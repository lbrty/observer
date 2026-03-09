import { resolve } from "node:path";

import { tanstackRouter } from "@tanstack/router-vite-plugin";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  if (mode !== "development" && mode !== "test") {
    const apiUrl = env.VITE_API_URL ?? "";
    if (!apiUrl.startsWith("https://")) {
      throw new Error(`VITE_API_URL must start with https:// for ${mode} builds. Got: "${apiUrl}"`);
    }
  }

  return {
    plugins: [
      react({
        babel: {
          plugins: [["babel-plugin-react-compiler"]],
        },
      }),
      tailwindcss(),
      tanstackRouter({
        routeFileIgnorePattern: ".+\\.test\\.tsx$",
      }),
    ],
    resolve: {
      alias: {
        "@": resolve(__dirname, "./src"),
      },
    },
  };
});
