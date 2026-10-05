// The gad IDE (@gad-lang/ide-vuetify) as a component of the admin's own Vue
// and Vuetify — window.Vue, window.Vuetify, not copies of them —: the editor
// of the files of a draft (admin/packages/gitedit), its API under base
// (base + "api/ide/…"). Loaded on demand by the admin's <vx-gad-ide>, which
// renders window.GadIde.component; it runs no code: the files and the
// language's operations only.
import { computed, defineComponent, h, onMounted, ref, watch, type PropType } from "vue";
import { useTheme } from "vuetify";
import { createHttpIdeApi, GadIde, type SerializedDockview, type Workspace } from "@gad-lang/ide-vuetify";

import "dockview-core/dist/styles/dockview.css";
import "./styles.css";

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

const GadIdeHost = defineComponent({
  name: "GadIdeHost",
  props: {
    // where the API is: base + "api/ide/…" ("/admin/site/site-files/ide/")
    base: { type: String, required: true },
    // the height it takes ("calc(100vh - 64px)"); its box's by default
    height: { type: String, default: "100%" },
    // the panels it has: no Output, no debugger — it runs no code
    panels: { type: Array as PropType<string[]>, default: () => ["explorer", "editor", "docs"] },
  },
  setup(props) {
    const api = createHttpIdeApi(props.base);
    // dark as the admin's theme is, and follows it when it changes
    const theme = useTheme();
    const dark = computed(() => theme.global.current.value.dark);
    const workspace = ref<Workspace>();
    const error = ref("");
    const layout = ref<SerializedDockview | null>(usableLayout(load("gadide.layout")));
    const config = ref<Record<string, unknown>>(load("gadide.config") ?? {});
    watch(layout, (v) => save("gadide.layout", v));
    watch(config, (v) => save("gadide.config", v), { deep: true });
    onMounted(async () => {
      try {
        workspace.value = await api.workspace();
      } catch (e) {
        error.value = String(e);
      }
    });
    return () =>
      h("div", { class: "gadide-host", style: { height: props.height } }, [
        error.value
          ? h("pre", { class: "gadide-error" }, error.value)
          : workspace.value &&
            h(GadIde, {
              api,
              workspace: workspace.value,
              dark: dark.value,
              runMode: "none",
              panels: props.panels,
              layoutConfig: layout.value,
              "onUpdate:layoutConfig": (v: SerializedDockview) => (layout.value = v),
              config: config.value,
              "onUpdate:config": (v: Record<string, unknown>) => (config.value = v),
            }),
      ]);
  },
});

declare const window: Window & { GadIde?: { component: typeof GadIdeHost } };
window.GadIde = { component: GadIdeHost };
