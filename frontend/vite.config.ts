import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { globSync } from "glob";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    manifest: true,
    rollupOptions: {
      input: [
        "src/main.css",
        "src/react/index.tsx",
        ...globSync(["src/vanilla/index.ts", "src/vanilla/**/*.ts"]),
      ],
    },
  },
});
