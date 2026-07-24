// bun test preload for the DOM tests (*.dom.test.ts):
//  1) a `.vue` loader — bun has no real SFC loader, so it would leave the
//     <template> uncompiled (the component then renders as its file path). We
//     compile each SFC (script setup + template) with @vue/compiler-sfc. The
//     corejs SFCs have no <style>, so styles are ignored.
//  2) happy-dom registered globally so the corejs Root can mount.
import { plugin } from "bun";
import { parse, compileScript, compileTemplate } from "@vue/compiler-sfc";
import { readFileSync } from "node:fs";
import { GlobalRegistrator } from "@happy-dom/global-registrator";

plugin({
  name: "vue-sfc",
  setup(build) {
    build.onLoad({ filter: /\.vue$/ }, (args) => {
      const source = readFileSync(args.path, "utf8");
      const id = args.path;
      const { descriptor } = parse(source, { filename: id });

      const scoped = descriptor.styles.some((s) => s.scoped);
      const scopeId = scoped ? "data-v-" + hash(id) : undefined;

      // <script setup> / <script> -> a component options object bound to `_sfc_main`.
      const script = compileScript(descriptor, {
        id,
        inlineTemplate: false,
        genDefaultAs: "_sfc_main",
      });

      const template = compileTemplate({
        source: descriptor.template?.content ?? "",
        filename: id,
        id,
        scoped,
        slotted: descriptor.slotted,
        compilerOptions: { scopeId, bindingMetadata: script.bindings },
      });

      const code =
        script.content +
        "\n" +
        template.code.replace("export function render", "function render") +
        "\n_sfc_main.render = render" +
        (scopeId ? `\n_sfc_main.__scopeId = ${JSON.stringify(scopeId)}` : "") +
        `\n_sfc_main.__file = ${JSON.stringify(id)}` +
        "\nexport default _sfc_main\n";

      return { contents: code, loader: "ts" };
    });
  },
});

function hash(s: string): string {
  let h = 0;
  for (let i = 0; i < s.length; i++) h = (Math.imul(31, h) + s.charCodeAt(i)) | 0;
  return (h >>> 0).toString(36);
}

GlobalRegistrator.register({ url: "http://localhost/" });

// allow the tests' fetch to reach the Go server on 127.0.0.1:<ephemeral> (the
// DOM's same-origin policy would otherwise block the cross-origin request).
const w = globalThis as any;
if (w.happyDOM?.settings?.fetch) {
  w.happyDOM.settings.fetch.disableSameOriginPolicy = true;
}
// happy-dom lacks ResizeObserver, which Vuetify components use.
if (typeof w.ResizeObserver === "undefined") {
  w.ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
}
