// DOM integration tests (bun + happy-dom) for the LISTING's form hosts.
//
// A listing renders ONE host per action, not one per row:
//   - `creating`      → the New button (see formhost.dom.test.ts);
//   - `itemDetailing` → clicking a row (when the model has a detailing);
//   - `itemEditing`   → the row menu's Edit (or the row click when there is no
//                       detailing).
//
// A row does not carry its own overlay plaid: it points the shared host at its
// record and turns it on — `<scope>.id = "<id>"; <scope>.show = true` — which
// mounts the guarded block, loads that record in ONE request, and destroys
// everything when turned off. Opening another row reuses the same host.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { defineComponent, inject, nextTick } from "vue";
import { flushPromises } from "@vue/test-utils";
import { installFetch, mountPresets, type ServerFetch } from "./dom";
import { startServer, type TestServer } from "./server";

let server: TestServer;
let sf: ServerFetch;
let fetchCount = 0;

beforeAll(async () => {
  server = await startServer("listeditor");
  sf = installFetch(server);
  const patched = globalThis.fetch;
  globalThis.fetch = ((input: any, init?: any) => {
    fetchCount++;
    return patched(input, init);
  }) as typeof fetch;
});

afterAll(() => {
  sf?.restore();
  server?.stop();
});

const Listing = defineComponent({
  template: `
    <div>
      <button id="loadBtn" @click='plaid().url("/admin/products").eventFunc("presets_OpenListingDialog").query("target_portal","page").go()'>load</button>
      <go-plaid-portal :visible="true" portal-name="page"></go-plaid-portal>
    </div>`,
  setup() {
    return { plaid: inject("plaid") };
  },
});

async function settle(ms = 700) {
  await flushPromises();
  await nextTick();
  await new Promise((r) => setTimeout(r, ms));
  await flushPromises();
  await nextTick();
}

function overlay(wrapper: any, title: string) {
  return wrapper
    .findAllComponents({ name: "vx-dialog" })
    .find((c: any) => String(c.props("title") ?? "").includes(title));
}

async function mountListing() {
  const wrapper = mountPresets(`<Listing></Listing>`, { Listing });
  await nextTick();
  await wrapper.find("#loadBtn").trigger("click");
  await settle(400);
  return wrapper;
}

describe("listing — per-item form host", () => {
  it("a row opens the shared host for its record (one request), and closing destroys it", async () => {
    const wrapper = await mountListing();

    // nothing is open yet: the row does not load anything until it flips the host
    expect(overlay(wrapper, "Product#1")).toBeFalsy();

    // click the row -> the shared host loads THAT record in a single request
    const before = fetchCount;
    const row = wrapper.findAll("td")[0];
    await row.trigger("click");
    await settle();
    expect(fetchCount - before).toBe(1);
    const detail = overlay(wrapper, "Product#1");
    expect(detail).toBeTruthy();

    // closing it destroys the overlay (the host's scope is its closer)
    detail!.vm.$emit("update:modelValue", false);
    await settle(200);
    expect(overlay(wrapper, "Product#1")).toBeFalsy();

    wrapper.unmount();
  }, 30000);

  it("the detailing opened from a row hosts its own edit form", async () => {
    const wrapper = await mountListing();

    const row = wrapper.findAll("td")[0];
    await row.trigger("click");
    await settle();
    expect(overlay(wrapper, "Product#1")).toBeTruthy();

    // the detailing rendered its own edit host: its Edit button only flips it,
    // and doing so loads the edit form in one more request
    const editBtn = wrapper.find('[data-event="edit"]');
    expect(editBtn.exists()).toBe(true);
    const before = fetchCount;
    await editBtn.trigger("click");
    await settle();
    expect(fetchCount - before).toBe(1);
    expect(overlay(wrapper, "Editing Product")).toBeTruthy();

    wrapper.unmount();
  }, 30000);
});
