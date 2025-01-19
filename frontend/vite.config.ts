import { svelte } from "@sveltejs/vite-plugin-svelte";
import react from "@vitejs/plugin-react";
import vue from "@vitejs/plugin-vue";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react(), vue(), svelte()],
  build: {
    manifest: true,
    rollupOptions: {
      input: [
        "src/main.ts",
        "src/main.css",
        "src/react/index.tsx",
        "src/vue/index.ts",
        "src/svelte/index.ts",
      ],
    },
  },
});
