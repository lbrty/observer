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
          manualChunks(id) {
            if (id.includes("node_modules")) {
              if (id.includes("react-dom") || id.includes("/react/")) return "vendor-react";
              if (id.includes("@tanstack/react-router") || id.includes("@tanstack/react-query") || id.includes("@tanstack/react-table")) return "vendor-router";
              if (id.includes("/d3") || id.includes("d3-sankey")) return "vendor-d3";
              if (id.includes("@base-ui") || id.includes("cmdk") || id.includes("react-day-picker")) return "vendor-ui";
              if (id.includes("i18next")) return "vendor-i18n";
              if (id.includes("@phosphor-icons")) return "vendor-icons";
            }
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
