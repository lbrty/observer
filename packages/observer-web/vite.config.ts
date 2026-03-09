import { resolve } from "node:path";

import { tanstackRouter } from "@tanstack/router-vite-plugin";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const backendUrl = process.env.VITE_API_URL || "http://localhost:9000";

export default defineConfig(() => {
  return {
    server: {
      proxy: {
        "/auth": backendUrl,
        "/admin": backendUrl,
        "/projects": backendUrl,
        "/my": backendUrl,
        "/health": backendUrl,
      },
    },
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
