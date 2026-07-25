// DOM integration tests (bun + happy-dom) for NESTED form hosts.
//
// A form host is not a singleton: several listings can be alive at the same
// time — the one on the page, a nested list opened in a dialog, and from THAT
// list a record's detail. Each one declares its own state as a slot variable
// (`<user-component :scope='{$presetsItemDetailing: …}'>`), so its guarded block
// is watched by that render alone.
//
// Were the state shared on the app-global `vars`, every live listing would
// declare the same well-known name and guard its block with `v-if` on it:
// opening a row in ONE of them would mount them ALL — several requests, and the
// overlay rendered into portals nobody is looking at. These tests pin that down.

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

// two independent listings, each loaded into its own portal — the shape of a
// page listing with a nested list dialog opened on top of it.
const TwoListings = defineComponent({
  template: `
    <div>
      <button id="loadA" @click='plaid().url("/admin/products").eventFunc("presets_OpenListingDialog").query("target_portal","pa").go()'>A</button>
      <button id="loadB" @click='plaid().url("/admin/products").eventFunc("presets_OpenListingDialog").query("target_portal","pb").go()'>B</button>
      <go-plaid-portal :visible="true" portal-name="pa"></go-plaid-portal>
      <go-plaid-portal :visible="true" portal-name="pb"></go-plaid-portal>
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

function dialogs(wrapper: any, title: string) {
  return wrapper
    .findAllComponents({ name: "vx-dialog" })
    .filter((c: any) => String(c.props("title") ?? "").includes(title));
}

async function mountBoth() {
  const wrapper = mountPresets(`<TwoListings></TwoListings>`, { TwoListings });
  await nextTick();
  await wrapper.find("#loadA").trigger("click");
  await settle(400);
  await wrapper.find("#loadB").trigger("click");
  await settle(400);
  return wrapper;
}

describe("nested form hosts", () => {
  it("two live listings own separate state: a row opens ONE overlay, once", async () => {
    const wrapper = await mountBoth();

    // both listings are up, each with its own rows
    const tables = wrapper.findAll("table");
    expect(tables.length).toBe(2);
    expect(dialogs(wrapper, "Product#1").length).toBe(0);

    // open a record from the SECOND listing
    const before = fetchCount;
    await tables[1].findAll("td")[0].trigger("click");
    await settle();

    // exactly one request and one overlay: the first listing's host did not
    // react (it would have, sharing the variable), so nothing loaded twice
    expect(fetchCount - before).toBe(1);
    expect(dialogs(wrapper, "Product#1").length).toBe(1);

    wrapper.unmount();
  }, 30000);

  it("closing the overlay of one listing leaves the other listing usable", async () => {
    const wrapper = await mountBoth();
    const tables = wrapper.findAll("table");

    await tables[1].findAll("td")[0].trigger("click");
    await settle();
    const openedFromB = dialogs(wrapper, "Product#1")[0];
    expect(openedFromB).toBeTruthy();

    // turning B's host off destroys B's overlay only
    openedFromB.vm.$emit("update:modelValue", false);
    await settle(200);
    expect(dialogs(wrapper, "Product#1").length).toBe(0);

    // the FIRST listing still opens its own — a fresh single request
    const before = fetchCount;
    await tables[0].findAll("td")[0].trigger("click");
    await settle();
    expect(fetchCount - before).toBe(1);
    expect(dialogs(wrapper, "Product#1").length).toBe(1);

    wrapper.unmount();
  }, 30000);

  it("the detail opened from a nested listing hosts its edit form on its own state", async () => {
    const wrapper = await mountBoth();
    const tables = wrapper.findAll("table");

    await tables[1].findAll("td")[0].trigger("click");
    await settle();

    // the detailing rendered inside the dialog carries its own edit host; the
    // button flips THAT state and loads the form once
    const editBtn = wrapper.find('[data-event="edit"]');
    expect(editBtn.exists()).toBe(true);
    const before = fetchCount;
    await editBtn.trigger("click");
    await settle();
    expect(fetchCount - before).toBe(1);
    expect(dialogs(wrapper, "Editing Product").length).toBe(1);

    wrapper.unmount();
  }, 30000);
});
