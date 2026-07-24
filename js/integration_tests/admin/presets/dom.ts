// DOM test helpers (bun + happy-dom): mount the real corejs Root with a light
// vx-dialog stub so the edit form's fields render, and point fetch at the Go
// server. Imported by the *.dom.test.ts specs. The corejs SFCs are compiled by
// the vue-sfc bun plugin in setup.ts.

import { defineComponent, inject } from "vue";
import { mount } from "@vue/test-utils";
import { createVuetify } from "vuetify";
import * as vuetifyComponents from "vuetify/components";
import * as vuetifyDirectives from "vuetify/directives";
// @ts-expect-error corejs @ alias resolved by corejs tsconfig
import { Root, plaidPlugin } from "../../../corejs/src/app";
import { type TestServer } from "./server";

const vuetify = createVuetify({
  components: vuetifyComponents,
  directives: vuetifyDirectives,
});

// The full vuetifyx library is too heavy to register here, but the edit form only
// needs `vx-dialog` to render its slots — the fields inside are plain
// vuetify/corejs components. So register a light stub.
export const VxDialogStub = defineComponent({
  name: "vx-dialog",
  props: {
    modelValue: { type: Boolean, default: true },
    title: { type: String, default: "" },
  },
  template: `
    <div class="vx-dialog" v-if="modelValue" :data-title="title">
      <div class="vx-dialog-toolbar"><slot name="appendToolbar"></slot></div>
      <div class="vx-dialog-body"><slot name="body"></slot></div>
      <slot></slot>
    </div>`,
});

export function mountPresets(template: string, components: Record<string, any> = {}) {
  return mount(Root, {
    props: { initialTemplate: template },
    global: {
      plugins: [plaidPlugin, vuetify],
      components: { "vx-dialog": VxDialogStub, ...components },
    },
  });
}

export interface ServerFetch {
  lastUpdate(): FormData | null;
  restore(): void;
}

// installFetch points corejs's fetch at the Go server and works around a DOM
// limitation: happy-dom's FormData is not serialized with empty fields, which
// would make the server ignore cleared inputs. So it re-encodes the body as
// url-encoded — the Go server accepts non-multipart submits (UnmarshalForm falls
// back to r.Form) — and captures the last Update body for assertions.
export function installFetch(server: TestServer): ServerFetch {
  const real = globalThis.fetch;
  let last: FormData | null = null;
  globalThis.fetch = (async (input: any, init?: any) => {
    let url = typeof input === "string" ? input : input.url;
    if (typeof url === "string" && url.startsWith("/")) url = server.url + url;
    if (init?.body && typeof (init.body as any).entries === "function") {
      const fd = init.body as FormData;
      if (typeof url === "string" && url.includes("presets_Update")) last = fd;
      const usp = new URLSearchParams();
      for (const [k, v] of fd.entries()) usp.append(k, String(v));
      init = { ...init, body: usp };
    }
    return real(url, init);
  }) as typeof fetch;
  return {
    lastUpdate: () => last,
    restore: () => {
      globalThis.fetch = real;
    },
  };
}

// makeEditButton builds a button that fires presets_EditForm into the outer
// portal (exposes plaid() to its template like the corejs specs). `url` selects
// the model listing.
export function makeEditButton(url: string) {
  return defineComponent({
    template: `
      <div>
        <button id="editBtn" @click='plaid()
          .url("${url}")
          .eventFunc("presets_EditForm")
          .query("id", "1")
          .query("overlay", "Dialog")
          .query("target_portal", "outerPortal")
          .go()'>Edit</button>
        <go-plaid-portal :visible="true" portal-name="outerPortal"></go-plaid-portal>
      </div>`,
    setup() {
      return { plaid: inject("plaid") };
    },
  });
}
