import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueJsx from "@vitejs/plugin-vue-jsx";
import { resolve } from "node:path";
import { existsSync } from "node:fs";

// The IDE of gad (web/ide-vuetify) and its plugins, from the source of a gad
// checkout: GAD_DIR, or the one beside this workspace. They are not on npm yet,
// and their source is not in gad's Go module (they are its git submodules).
const gad = resolve(process.env.GAD_DIR || resolve(__dirname, "../../../../../../../gade/src/github.com/gad-lang/gad"));
const src = (rel: string) => {
  const p = resolve(gad, rel);
  if (!existsSync(p)) throw new Error(`gadide: no ${p} (set GAD_DIR to a gad checkout with its submodules)`);
  return p;
};

// A library for the admin's page: Vue and Vuetify are the page's globals
// (window.Vue, window.Vuetify) — one of each, the admin's theme and its
// components —, the rest in gadide.js and gadide.css.
export default defineConfig({
  plugins: [vue(), vueJsx()],
  define: { "process.env.NODE_ENV": JSON.stringify("production") },
  resolve: {
    alias: {
      "@gad-lang/ide-vuetify": src("web/ide-vuetify/src/index.ts"),
      "@gad-lang/codemirror-gad": src("web/plugins/js/codemirror-gad/src/index.ts"),
      "@gad-lang/prism-gad": src("web/plugins/js/prism-gad/src/index.ts"),
    },
    // one of each: CodeMirror's state fields need it
    dedupe: ["@codemirror/state", "@codemirror/view", "@codemirror/language",
      "@codemirror/autocomplete", "@codemirror/lint", "@codemirror/commands",
      "@lezer/common", "@lezer/highlight", "@lezer/lr", "dockview-vue", "dockview-core", "prismjs"],
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
    chunkSizeWarningLimit: 4096,
    cssCodeSplit: false,
    lib: { entry: resolve(__dirname, "src/main.ts"), formats: ["iife"], name: "gadidejs", fileName: () => "gadide.js" },
    rollupOptions: {
      external: ["vue", "vuetify", "vuetify/components", "vuetify/directives"],
      output: {
        globals: {
          vue: "Vue",
          vuetify: "Vuetify",
          "vuetify/components": "Vuetify.components",
          "vuetify/directives": "Vuetify.directives",
        },
        assetFileNames: () => "gadide.css",
      },
    },
  },
});
