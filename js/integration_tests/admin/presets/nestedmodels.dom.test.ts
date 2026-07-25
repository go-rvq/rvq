// DOM integration tests (bun + happy-dom) for TWO LEVELS of nested models —
// has-many fields taken as child models with their own listing/detailing
// (helper.NewNestedSliceBuilder), the shape an application reaches with
// `Root → Root.Children → Root.Children.GrandChildren`:
//
//   Survey#1
//     └─ Places listing (dialog)               ← level 1
//          ├─ Place#N detail (dialog)
//          └─ Products listing of a place      ← level 2
//               └─ Product#N detail (dialog)
//
// Every level renders the same four kinds of host. Were their state shared, a
// row of one level would drive another level's overlay — several requests and
// duplicated overlays. Each level declares its own slot variables instead, and
// these tests keep it that way.
//
// Level 2 is reached the way a user reaches it: every row of the level-1 listing
// carries a "…" menu whose entry for the nested field opens THAT record's
// listing.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { defineComponent, inject, nextTick } from "vue";
import { flushPromises } from "@vue/test-utils";
import { installFetch, mountPresets, type ServerFetch } from "./dom";
import { startServer, type TestServer } from "./server";

let server: TestServer;
let sf: ServerFetch;
let fetchCount = 0;

beforeAll(async () => {
  server = await startServer("nestedmodels");
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

const Root = defineComponent({
  template: `
    <div>
      <button id="places" @click='plaid().url("/admin/surveys/1/places").eventFunc("presets_OpenListingDialog").query("target_portal","l1").query("overlay","Dialog").go()'>places</button>
      <go-plaid-portal :visible="true" portal-name="l1"></go-plaid-portal>
    </div>`,
  setup() {
    return { plaid: inject("plaid") };
  },
});

async function settle(ms = 900) {
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

// Dialogs STACK: a nested listing opened from a row is rendered inside the DOM
// subtree of the dialog it came from, so "the table of dialog X" is ambiguous —
// tables are picked by their content instead.
function tableWith(wrapper: any, needle: string) {
  return wrapper.findAll("table").find((t: any) => t.text().includes(needle));
}

// clicks the "…" menu of the first row and picks the entry with this label. The
// menu is rendered by Vuetify into an overlay outside the wrapper, so its items
// are looked up on the document.
async function rowMenu(table: any, label: string) {
  await table.findAll("tbody button")[0].trigger("click");
  await settle(400);
  const item = Array.from(document.body.querySelectorAll(".v-list-item")).find(
    (n) => (n.textContent ?? "").trim() === label,
  );
  expect(item).toBeTruthy();
  (item as HTMLElement).click();
  await settle();
}

async function openBothLevels() {
  const wrapper = mountPresets(`<Root></Root>`, { Root });
  await nextTick();
  await wrapper.find("#places").trigger("click");
  await settle();

  // level 2: the first row's "…" menu → the nested listing of THAT place
  await rowMenu(tableWith(wrapper, "Place "), "Products");
  return wrapper;
}

describe("two levels of nested models", () => {
  it("both levels render their own listing, each with its own rows", async () => {
    const wrapper = await openBothLevels();

    // both dialogs are stacked, each with its own table
    expect(dialogs(wrapper, "Listing Places").length).toBe(1);
    expect(dialogs(wrapper, "Listing Products").length).toBe(1);

    const products = tableWith(wrapper, "-one");
    expect(products).toBeTruthy();

    // the level-2 listing is scoped to the place its row menu was opened from
    const rows = products!.findAll("tbody tr").length;
    expect(rows).toBeGreaterThan(0);
    expect(rows).toBeLessThan(3);

    wrapper.unmount();
  }, 40000);

  it("a row of the LEVEL-2 listing opens one overlay, once", async () => {
    const wrapper = await openBothLevels();

    const openBefore = wrapper.findAllComponents({ name: "vx-dialog" }).length;
    const before = fetchCount;
    await tableWith(wrapper, "-one")!.findAll("td")[0].trigger("click");
    await settle();

    // exactly one request and one new overlay: the level-1 listing's host did
    // not react to it (sharing the variable, it would have opened a place too)
    expect(fetchCount - before).toBe(1);
    expect(wrapper.findAllComponents({ name: "vx-dialog" }).length).toBe(openBefore + 1);

    wrapper.unmount();
  }, 40000);

  it("a row of the LEVEL-1 listing opens one overlay, once", async () => {
    const wrapper = await openBothLevels();

    const openBefore = wrapper.findAllComponents({ name: "vx-dialog" }).length;
    const before = fetchCount;
    await tableWith(wrapper, "Place ")!.findAll("td")[0].trigger("click");
    await settle();

    expect(fetchCount - before).toBe(1);
    expect(wrapper.findAllComponents({ name: "vx-dialog" }).length).toBe(openBefore + 1);

    wrapper.unmount();
  }, 40000);

  it("both details stack and each Edit button loads only its own form", async () => {
    const wrapper = await openBothLevels();

    // a place (from level 1) and a product (from level 2), both open at once —
    // dialogs stack, they do not replace each other
    await tableWith(wrapper, "Place ")!.findAll("td")[0].trigger("click");
    await settle();
    await tableWith(wrapper, "-one")!.findAll("td")[0].trigger("click");
    await settle();

    expect(dialogs(wrapper, "Listing Places").length).toBe(1);
    expect(dialogs(wrapper, "Listing Products").length).toBe(1);
    expect(dialogs(wrapper, "Nmplace").length).toBe(1);
    expect(dialogs(wrapper, "Nmproduct").length).toBe(1);

    // one Edit button per open detail, each bound to its own edit host
    const editBtns = wrapper.findAll('[data-event="edit"]');
    expect(editBtns.length).toBe(2);

    for (const btn of editBtns) {
      const before = fetchCount;
      await btn.trigger("click");
      await settle();
      // only that detail's form loaded — a shared variable would have loaded both
      expect(fetchCount - before).toBe(1);
    }

    wrapper.unmount();
  }, 40000);
});
