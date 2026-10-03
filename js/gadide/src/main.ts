// The gad IDE on the editor's draft: its API under this page (api/ide/…,
// relative), as served by admin/packages/gitedit — the files and the
// language's operations, no running of code.
import { createApp, defineComponent, h, onMounted, ref, watch } from "vue";
import { createVuetify } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";
import { aliases, mdi } from "vuetify/iconsets/mdi";
import { GadIde, httpIdeApi, type SerializedDockview, type Workspace } from "@gad-lang/ide-vuetify";

import "vuetify/styles";
import "@mdi/font/css/materialdesignicons.css";
import "dockview-core/dist/styles/dockview.css";
import "./styles.css";

const dark = window.matchMedia?.("(prefers-color-scheme: dark)").matches ?? false;

const load = <T,>(key: string): T | null => {
  try {
    const v = localStorage.getItem(key);
    return v ? (JSON.parse(v) as T) : null;
  } catch {
    return null;
  }
};
const save = (key: string, v: unknown) => {
  try {
    localStorage.setItem(key, JSON.stringify(v));
  } catch {
    /* private window */
  }
};

// a layout saved that shows neither the explorer nor the editor — closed, or
// of an older version — would open a blank IDE: the default one instead
const usableLayout = (l: SerializedDockview | null): SerializedDockview | null => {
  const panels = (l as { panels?: Record<string, unknown> } | null)?.panels;
  return panels && ("explorer" in panels || "editor" in panels) ? l : null;
};

// an error shows on the page, not a blank one
const showError = (msg: string) => {
  let el = document.getElementById("gadide-error");
  if (!el) {
    el = document.createElement("pre");
    el.id = "gadide-error";
    el.className = "gadide-error";
    document.body.prepend(el);
  }
  el.textContent += msg + "\n";
};
window.addEventListener("error", (e) => showError(String(e.error?.stack || e.message)));
window.addEventListener("unhandledrejection", (e) => showError(String(e.reason?.stack || e.reason)));

const App = defineComponent(() => {
  const workspace = ref<Workspace>();
  const error = ref("");
  const layout = ref<SerializedDockview | null>(usableLayout(load("gadide.layout")));
  const config = ref<Record<string, unknown>>(load("gadide.config") ?? {});
  watch(layout, (v) => save("gadide.layout", v));
  watch(config, (v) => save("gadide.config", v), { deep: true });
  onMounted(async () => {
    try {
      workspace.value = await httpIdeApi.workspace();
    } catch (e) {
      error.value = String(e);
    }
  });
  return () =>
    error.value
      ? h("pre", { class: "gadide-error" }, error.value)
      : workspace.value &&
        h(GadIde, {
          api: httpIdeApi,
          workspace: workspace.value,
          dark,
          runMode: "none",
          layoutConfig: layout.value,
          "onUpdate:layoutConfig": (v: SerializedDockview) => (layout.value = v),
          config: config.value,
          "onUpdate:config": (v: Record<string, unknown>) => (config.value = v),
        });
});

createApp(App)
  .use(createVuetify({ components, directives, icons: { defaultSet: "mdi", aliases, sets: { mdi } }, theme: { defaultTheme: dark ? "dark" : "light" } }))
  .mount("#app");
