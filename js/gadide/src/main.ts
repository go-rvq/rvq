// The gad IDE (@gad-lang/ide-vuetify) as a component of the admin's own Vue
// and Vuetify — window.Vue, window.Vuetify, not copies of them —: the editor
// of the files of a draft (admin/packages/gitedit), its API under base
// (base + "api/ide/…"). Loaded on demand by the admin's <vx-gad-ide>, which
// renders window.GadIde.component; it runs no code: the files and the
// language's operations only.
import { computed, defineComponent, h, onBeforeUnmount, onMounted, ref, resolveComponent, watch, type PropType } from "vue";
import { useTheme } from "vuetify";
import { useI18n } from "vue-i18n";
import { createHttpIdeApi, GadIde, type ExtraPanel, type HeaderAction, type IdeMessages, type SerializedDockview, type Workspace } from "@gad-lang/ide-vuetify";

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

// a layout saved that shows none of the panels of the mode (the IDE: the
// explorer, the editor; git: Changes, Git) — closed, or of an older version
// — would open a blank IDE: the default one instead
const usableLayout = (l: SerializedDockview | null, want: string[]): SerializedDockview | null => {
  const panels = (l as { panels?: Record<string, unknown> } | null)?.panels;
  return panels && want.some((id) => id in panels) ? l : null;
};

// HostPanel is a panel of the page in the IDE (extraPanels): a component of
// the admin's Vue app, by its name (tag), given props — placed below by
// default, its button in the IDE's bottom strip, its expand in its tab.
export interface HostPanel {
  id: string;
  label: string;
  icon?: string;
  tag: string;
  props?: Record<string, unknown>;
  placement?: "bottom" | "right" | "left";
}

// hostPanel is the IDE's ExtraPanel of p: its component the app's tag.
const hostPanel = (p: HostPanel): ExtraPanel => ({
  id: p.id,
  label: p.label,
  icon: p.icon,
  placement: p.placement,
  toolbarButton: true,
  tabExpand: true,
  component: defineComponent({
    name: "HostPanel_" + p.id.replace(/\W/g, "_"),
    setup() {
      return () => h("div", { class: "gadide-host-panel" }, [h(resolveComponent(p.tag) as never, p.props)]);
    },
  }),
});

const GadIdeHost = defineComponent({
  name: "GadIdeHost",
  props: {
    // where the API is: base + "api/ide/…" ("/admin/site/site-files/ide/")
    base: { type: String, required: true },
    // the height it takes ("calc(100vh - 64px)"); its box's by default;
    // "fit": from where it is to the bottom of the window (the page above it
    // whatever it is), again when the window or the page changes
    height: { type: String, default: "100%" },
    // the panels it has: no Output, no debugger — it runs no code
    panels: { type: Array as PropType<string[]>, default: () => ["explorer", "editor", "structure", "docs"] },
    // the texts of the IDE and of its Changes and Git panels, translated (the
    // admin's, gitedit's IdeMessages): merged into the admin's vue-i18n
    messages: { type: Object as PropType<Partial<IdeMessages>>, default: undefined },
    // buttons of the page in the IDE's header, at its right (HeaderAction)
    headerActions: { type: Array as PropType<HeaderAction[]>, default: () => [] },
    // panels of the page (HostPanel): a component of the app each
    extraPanels: { type: Array as PropType<HostPanel[]>, default: () => [] },
    // "git": the IDE of a repository — its Changes and Git panels alone,
    // open side by side, its layout kept apart from the IDE's
    mode: { type: String as PropType<"ide" | "git">, default: "ide" },
    // the Settings' tabs shown (panels, formatter, transpile, template); none: all
    settingsTabs: { type: Array as PropType<string[]>, default: () => [] },
  },
  setup(props) {
    const api = createHttpIdeApi(props.base);
    // the admin's words into its vue-i18n (window.VueI18n, the page's
    // locale): the IDE reads them there (gadIde, gadGit) — given as a prop
    // only with no vue-i18n
    let viaI18n = false;
    try {
      const i18n = useI18n({ useScope: "global" });
      if (props.messages) i18n.mergeLocaleMessage(String(i18n.locale.value), { gadIde: props.messages, gadGit: props.messages });
      viaI18n = true;
    } catch {
      /* no vue-i18n installed */
    }
    // dark as the admin's theme is, and follows it when it changes
    const theme = useTheme();
    const dark = computed(() => theme.global.current.value.dark);
    // "fit": the window's height less the IDE's top in the page
    const host = ref<HTMLElement>();
    const fitted = ref("");
    const fit = () => {
      const el = host.value;
      if (!el || props.height !== "fit") return;
      const top = el.getBoundingClientRect().top + window.scrollY;
      fitted.value = Math.max(320, window.innerHeight - top - 8) + "px";
    };
    let ro: ResizeObserver | undefined;
    onMounted(() => {
      if (props.height !== "fit") return;
      fit();
      window.addEventListener("resize", fit);
      // what is above it may grow or shrink (a bar, a message)
      ro = new ResizeObserver(fit);
      ro.observe(document.body);
    });
    onBeforeUnmount(() => {
      window.removeEventListener("resize", fit);
      ro?.disconnect();
    });
    const height = computed(() => (props.height === "fit" ? fitted.value || "calc(100vh - 64px)" : props.height));
    // the page's panels, made once (the IDE takes them when it starts)
    const hostPanels = props.extraPanels.map(hostPanel);
    const workspace = ref<Workspace>();
    const error = ref("");
    const git = props.mode === "git";
    const key = git ? "gadide.git." : "gadide.";
    const layout = ref<SerializedDockview | null>(
      usableLayout(load(key + "layout"), git ? ["changes", "git"] : ["explorer", "editor"]),
    );
    const config = ref<Record<string, unknown>>(load(key + "config") ?? {});
    watch(layout, (v) => save(key + "layout", v));
    watch(config, (v) => save(key + "config", v), { deep: true });
    onMounted(async () => {
      try {
        workspace.value = await api.workspace();
      } catch (e) {
        error.value = String(e);
      }
    });
    return () =>
      h("div", { ref: host, class: "gadide-host", style: { height: height.value } }, [
        error.value
          ? h("pre", { class: "gadide-error" }, error.value)
          : workspace.value &&
            h(GadIde, {
              api,
              workspace: workspace.value,
              dark: dark.value,
              runMode: "none",
              panels: git ? [] : props.panels,
              openPanels: git ? ["changes", "git"] : [],
              settingsTabs: props.settingsTabs,
              messages: viaI18n ? undefined : props.messages,
              headerActions: props.headerActions,
              extraPanels: hostPanels,
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
