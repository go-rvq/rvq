// DOM integration tests (bun + happy-dom) for the POST-SAVE REFRESH.
//
// After a successful save, what shows the record must show the change — without
// the page being reloaded. The saved form carries nothing about its opener: it
// runs the hooks in `onSaveCallbacks`, the list every form host appended to on
// the way down.
//
//   listing → NEW              → the listing shows the new row
//   listing → DETAIL → EDIT    → the detail shows the new data AND its listing
//   listing without detailing  → the row's EDIT still refreshes the listing
//   singleton with detailing   → EDIT refreshes the detail

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { defineComponent, inject, nextTick } from "vue";
import { flushPromises } from "@vue/test-utils";
import { installFetch, mountPresets, type ServerFetch } from "./dom";
import { startServer, type TestServer } from "./server";

let server: TestServer;
let sf: ServerFetch;
let reloads = 0;

beforeAll(async () => {
  server = await startServer("refresh");
  sf = installFetch(server);
  // nothing in these flows may ask the browser to reload the page
  (globalThis as any).location = { ...(globalThis as any).location, reload: () => reloads++ };
});

afterAll(() => {
  sf?.restore();
  server?.stop();
});

function listingAt(url: string) {
  return defineComponent({
    template: `
      <div>
        <button id="load" @click='plaid().url("${url}").eventFunc("presets_OpenListingDialog").query("target_portal","page").query("overlay","Dialog").go()'>load</button>
        <go-plaid-portal :visible="true" portal-name="page"></go-plaid-portal>
      </div>`,
    setup() {
      return { plaid: inject("plaid") };
    },
  });
}

async function settle(ms = 900) {
  await flushPromises();
  await nextTick();
  await new Promise((r) => setTimeout(r, ms));
  await flushPromises();
  await nextTick();
}

async function mountListing(url: string) {
  const Listing = listingAt(url);
  const wrapper = mountPresets(`<Listing></Listing>`, { Listing });
  await nextTick();
  await wrapper.find("#load").trigger("click");
  await settle();
  return wrapper;
}

// Types a value into the form's FIRST field and saves it. Vuetify fields carry
// no name attribute (they get generated ids), so the field is taken by position
// inside the innermost dialog — the form that was just opened.
async function fillAndSave(wrapper: any, value: string) {
  const dialogs = wrapper.findAllComponents({ name: "vx-dialog" });
  const form = dialogs[dialogs.length - 1];
  const input = form.findAll("input")[0];
  expect(input).toBeTruthy();
  await input.setValue(value);
  await settle(200);

  const save = [
    ...wrapper.findAll('[data-event="presets_Update"]'),
    ...wrapper.findAll('[data-event="presets_Create"]'),
  ];
  expect(save.length).toBeGreaterThan(0);
  await save[save.length - 1].trigger("click");
  await settle();
}

function tableText(wrapper: any) {
  return wrapper.findAll("table").map((t: any) => t.text()).join(" | ");
}

describe("NEW from a listing", () => {
  it("the listing shows the created record, with no page reload", async () => {
    const wrapper = await mountListing("/admin/articles");
    expect(tableText(wrapper)).not.toContain("Brand new");

    await wrapper.find('[data-event="new"]').trigger("click");
    await settle();
    await fillAndSave(wrapper, "Brand new");

    expect(tableText(wrapper)).toContain("Brand new");
    expect(reloads).toBe(0);

    wrapper.unmount();
  }, 40000);
});

describe("EDIT from a DETAIL opened by a listing", () => {
  it("refreshes the detail AND the listing that opened it", async () => {
    const wrapper = await mountListing("/admin/articles");

    // open a record's detail from its row
    await wrapper.findAll("table")[0].findAll("td")[0].trigger("click");
    await settle();
    const detail = wrapper
      .findAllComponents({ name: "vx-dialog" })
      .find((d: any) => !String(d.props("title") ?? "").startsWith("Listing"));
    expect(detail).toBeTruthy();

    // edit it from the detail
    await wrapper.find('[data-event="edit"]').trigger("click");
    await settle();
    await fillAndSave(wrapper, "Detail edited");

    // the detail body shows the new value…
    const detailNow = wrapper
      .findAllComponents({ name: "vx-dialog" })
      .find((d: any) => !String(d.props("title") ?? "").startsWith("Listing"));
    expect(detailNow!.text()).toContain("Detail edited");
    // …and so does the listing underneath
    expect(tableText(wrapper)).toContain("Detail edited");
    expect(reloads).toBe(0);

    wrapper.unmount();
  }, 40000);
});

describe("a listing whose model has no detailing", () => {
  it("still refreshes when the row's EDIT is saved", async () => {
    const wrapper = await mountListing("/admin/tags");

    // with no detailing, the row opens the edit form itself
    await wrapper.findAll("table")[0].findAll("td")[0].trigger("click");
    await settle();
    await fillAndSave(wrapper, "Tag edited");

    expect(tableText(wrapper)).toContain("Tag edited");
    expect(reloads).toBe(0);

    wrapper.unmount();
  }, 40000);
});

describe("singleton with a detailing", () => {
  it("saving its EDIT refreshes the DETAIL", async () => {
    const Single = defineComponent({
      template: `
        <div>
          <button id="load" @click='plaid().url("/admin/site").eventFunc("presets_Detailing").query("target_portal","page").query("overlay","Dialog").go()'>load</button>
          <go-plaid-portal :visible="true" portal-name="page"></go-plaid-portal>
        </div>`,
      setup() {
        return { plaid: inject("plaid") };
      },
    });

    const wrapper = mountPresets(`<Single></Single>`, { Single });
    await nextTick();
    await wrapper.find("#load").trigger("click");
    await settle();

    await wrapper.find('[data-event="edit"]').trigger("click");
    await settle();
    await fillAndSave(wrapper, "Singleton edited");

    const detail = wrapper.findAllComponents({ name: "vx-dialog" })[0];
    expect(detail.text()).toContain("Singleton edited");
    expect(reloads).toBe(0);

    wrapper.unmount();
  }, 40000);
});
