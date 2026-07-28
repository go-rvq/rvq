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
  emits: ["update:modelValue"],
  // the close button stands for the real dialog's chrome: closing writes `false`
  // through `v-model='closer.show'`, which is what turns the overlay off.
  template: `
    <div class="vx-dialog" v-if="modelValue" :data-title="title">
      <button class="vx-dialog-close" @click="$emit('update:modelValue', false)"></button>
      <div class="vx-dialog-toolbar"><slot name="appendToolbar"></slot></div>
      <div class="vx-dialog-body"><slot name="body"></slot></div>
      <slot></slot>
    </div>`,
});

// Como o vx-dialog: o suficiente para o teste ver o overlay abrir e fechar.
export const VxNavigationDrawerStub = defineComponent({
  name: "vx-navigation-drawer",
  props: {
    modelValue: { type: Boolean, default: false },
    title: { type: String, default: "" },
  },
  emits: ["update:modelValue"],
  template: `
    <div class="vx-drawer" v-if="modelValue" :data-title="title">
      <button class="vx-drawer-close" @click="$emit('update:modelValue', false)"></button>
      <div class="vx-drawer-toolbar"><slot name="appendToolbar"></slot></div>
      <div class="vx-drawer-body"><slot name="body"></slot></div>
      <slot></slot>
    </div>`,
});

export function mountPresets(template: string, components: Record<string, any> = {}) {
  return mount(Root, {
    props: { initialTemplate: template },
    global: {
      plugins: [plaidPlugin, vuetify],
      components: {
        "vx-dialog": VxDialogStub,
        "vx-navigation-drawer": VxNavigationDrawerStub,
        ...components,
      },
      // the real app declares these as global properties (see createWebApp in
      // corejs/src/app.ts); without them a rendered portal that references
      // `presetsListing` cannot resolve it and its content fails to render.
      config: {
        globalProperties: {
          onSaveCallbacks: [],
          presetsListing: null,
          presetsDetailing: null,
          presetsCreating: null,
          presetsEditing: null,
        },
      },
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
