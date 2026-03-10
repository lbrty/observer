import { resolve } from "node:path";

import { tanstackRouter } from "@tanstack/router-vite-plugin";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig(() => {
  return {
    build: {
      outDir: "../../internal/spa/dist",
      emptyOutDir: true,
      rollupOptions: {
        output: {
          manualChunks: {
            "vendor-react":   ["react", "react-dom"],
            "vendor-router":  ["@tanstack/react-router", "@tanstack/react-query", "@tanstack/react-table"],
            "vendor-d3":      ["d3", "d3-sankey"],
            "vendor-ui":      ["@base-ui/react", "cmdk", "react-day-picker"],
            "vendor-i18n":    ["i18next", "react-i18next"],
            "vendor-icons":   ["@phosphor-icons/react"],
          },
        },
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
