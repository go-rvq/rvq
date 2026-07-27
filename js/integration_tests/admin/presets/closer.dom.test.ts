// DOM integration tests (bun + happy-dom) for the ADDRESS BAR of overlays.
//
// A dialog/drawer that shows a listing, a record's detail, its edit form or a
// create form is the same thing a PAGE shows — so while it is open the address
// bar carries that page's clean URL (`/admin/articles/1`, not an event), and
// closing it puts back the address that was there before. Overlays stack, so
// the addresses unwind LIFO.
//
// Drives the real corejs runtime against the Go fixture server: the URLs come
// from the form hosts the listing/detailing render (presets.FormHostBuilder.URL).

import { afterAll, beforeAll, beforeEach, describe, expect, it } from "bun:test";
import { defineComponent, inject, nextTick } from "vue";
import { flushPromises } from "@vue/test-utils";
import { installFetch, mountPresets, type ServerFetch } from "./dom";
import { startServer, type TestServer } from "./server";
// @ts-expect-error corejs @ alias resolved by corejs tsconfig
import { closerUrlStack, resetCloserUrlStack } from "../../../corejs/src/closer-url";

let server: TestServer;
let sf: ServerFetch;

beforeAll(async () => {
  server = await startServer("refresh");
  sf = installFetch(server);
});

afterAll(() => {
  sf?.restore();
  server?.stop();
});

beforeEach(() => {
  resetCloserUrlStack();
  window.history.replaceState(null, "", "/admin/articles");
});

async function settle(ms = 900) {
  await flushPromises();
  await nextTick();
  await new Promise((r) => setTimeout(r, ms));
  await flushPromises();
  await nextTick();
}

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

async function mountListing(url: string) {
  const Listing = listingAt(url);
  const wrapper = mountPresets(`<Listing></Listing>`, { Listing });
  await nextTick();
  await wrapper.find("#load").trigger("click");
  await settle();
  return wrapper;
}

function dialogs(wrapper: any) {
  return wrapper.findAllComponents({ name: "vx-dialog" });
}

// closes the innermost overlay the way the real dialog does (v-model → closer)
async function closeTop(wrapper: any) {
  const all = dialogs(wrapper);
  await all[all.length - 1].find(".vx-dialog-close").trigger("click");
  await settle(200);
}

function path() {
  return window.location.pathname;
}

describe("a record opened from a listing", () => {
  it("shows the detail page's address, then the edit page's, and unwinds LIFO", async () => {
    const wrapper = await mountListing("/admin/articles");
    // the listing itself was opened by the test, not by a host: the address is
    // still the page the test started on
    expect(path()).toBe("/admin/articles");

    // the row opens the record's DETAIL: its page is /admin/articles/{id}
    await wrapper.findAll("table")[0].findAll("td")[0].trigger("click");
    await settle();
    const detailPath = path();
    expect(detailPath).toMatch(/^\/admin\/articles\/\d+$/);

    // and its EDIT, from inside the detail: /admin/articles/{id}/edit
    await wrapper.find('[data-event="edit"]').trigger("click");
    await settle();
    expect(path()).toBe(detailPath + "/edit");
    expect(closerUrlStack().map((e: any) => e.url)).toEqual([
      detailPath,
      detailPath + "/edit",
    ]);

    // closing gives each address back, innermost first
    await closeTop(wrapper);
    expect(path()).toBe(detailPath);

    await closeTop(wrapper);
    expect(path()).toBe("/admin/articles");
    expect(closerUrlStack().length).toBe(0);

    wrapper.unmount();
  }, 40000);
});

describe("the create form", () => {
  it("shows the /new page address while it is open", async () => {
    const wrapper = await mountListing("/admin/articles");

    await wrapper.find('[data-event="new"]').trigger("click");
    await settle();
    expect(path()).toBe("/admin/articles/new");

    await closeTop(wrapper);
    expect(path()).toBe("/admin/articles");

    wrapper.unmount();
  }, 40000);
});

describe("a model with no detailing", () => {
  it("the row opens the EDIT page's address", async () => {
    const wrapper = await mountListing("/admin/tags");

    await wrapper.findAll("table")[0].findAll("td")[0].trigger("click");
    await settle();
    expect(path()).toMatch(/^\/admin\/tags\/\d+\/edit$/);

    await closeTop(wrapper);
    expect(path()).toBe("/admin/articles");

    wrapper.unmount();
  }, 40000);
});

describe("the addresses are real pages", () => {
  it("what the overlay showed is what the address serves", async () => {
    const first = await (await fetch(server.url + "/admin/articles?__execute_event__=__reload__")).text();
    const id = (first.match(/\/admin\/articles\/(\d+)/) || [])[1] || "1";
    for (const url of [`/admin/articles/${id}`, `/admin/articles/${id}/edit`, "/admin/articles/new"]) {
      const res = await fetch(server.url + url);
      expect(`${url} -> ${res.status}`).toBe(`${url} -> 200`);
    }
  }, 20000);
});
