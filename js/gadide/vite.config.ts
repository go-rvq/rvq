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

export default defineConfig({
  // served under the page of the admin: its URLs relative
  base: "./",
  plugins: [vue(), vueJsx()],
  resolve: {
    alias: {
      "@gad-lang/ide-vuetify": src("web/ide-vuetify/src/index.ts"),
      "@gad-lang/codemirror-gad": src("web/plugins/js/codemirror-gad/src/index.ts"),
      "@gad-lang/prism-gad": src("web/plugins/js/prism-gad/src/index.ts"),
    },
    // one of each: CodeMirror's state fields and Vuetify's injection need it
    dedupe: ["vue", "vuetify", "@codemirror/state", "@codemirror/view", "@codemirror/language",
      "@codemirror/autocomplete", "@codemirror/lint", "@codemirror/commands",
      "@lezer/common", "@lezer/highlight", "@lezer/lr", "dockview-vue", "dockview-core", "prismjs"],
  },
  server: { fs: { allow: [__dirname, gad] } },
  build: { outDir: "dist", emptyOutDir: true, chunkSizeWarningLimit: 4096 },
});
