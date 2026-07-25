// DOM integration tests (bun + happy-dom) for the FORM HOST: the page owns the
// form's `closer` in a scope var, so a button only has to flip it.
//
//   EDIT — the detailing hosts the edit form in `editing`:
//     click Edit  → `editing.show = true`  → the guarded block mounts, its
//                    run-script loads the Edit form into the host portal (ONE
//                    request — no wrapper round-trip) → the dialog appears;
//     `editing.show = false` → the block unmounts → form and its scope are
//                    destroyed completely;
//     `editing.show = true` again → a FRESH form is loaded.
//
//   NEW — the listing hosts the create form in `creating`, same semantics.
//
// Because the host owns the closer (the portal seeds `closer` in the content
// scope and ParamCloserProvided keeps the responder from creating a child one),
// closing the dialog from inside also destroys the form.

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
  // count the requests, to prove opening the form costs a single round-trip
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

// Loader mounts a page (detailing or listing) into a portal, then the tests
// interact with the rendered host exactly like a user would.
function makeLoader(event: string, query: Record<string, string>) {
  const q = Object.entries(query)
    .map(([k, v]) => `.query(${JSON.stringify(k)}, ${JSON.stringify(v)})`)
    .join("");
  return defineComponent({
    template: `
      <div>
        <button id="loadBtn" @click='plaid().url("/admin/products").eventFunc(${JSON.stringify(event)})${q}.query("target_portal","page").go()'>load</button>
        <go-plaid-portal :visible="true" portal-name="page"></go-plaid-portal>
      </div>`,
    setup() {
      return { plaid: inject("plaid") };
    },
  });
}

async function settle(ms = 700) {
  await flushPromises();
  await nextTick();
  await new Promise((r) => setTimeout(r, ms));
  await flushPromises();
  await nextTick();
}

// mountPage loads the given page event into the wrapper and returns it.
async function mountPage(event: string, query: Record<string, string>) {
  const Page = makeLoader(event, query);
  const wrapper = mountPresets(`<Page></Page>`, { Page });
  await nextTick();
  await wrapper.find("#loadBtn").trigger("click");
  await settle(300);
  return wrapper;
}

// formDialog finds the HOSTED form overlay by its title — the page itself is an
// overlay too, so we must not match that one.
function formDialog(wrapper: any, title: string) {
  return wrapper
    .findAllComponents({ name: "vx-dialog" })
    .find((c: any) => String(c.props("title") ?? "").includes(title));
}

// formOpen reports whether the hosted form is currently mounted (the host's
// v-if unmounts it entirely when its scope var is turned off).
function formOpen(wrapper: any, title: string): boolean {
  return !!formDialog(wrapper, title);
}

// closeForm does what closing the overlay does: turns the host's closer off.
async function closeForm(wrapper: any, title: string) {
  const d = formDialog(wrapper, title);
  expect(d).toBeTruthy();
  d!.vm.$emit("update:modelValue", false);
}

describe("form host — EDIT (detailing)", () => {
  it("the Edit button only flips editing.show, and that opens/destroys the form", async () => {
    // load the detailing as a page (not an overlay), so the only overlay around
    // is the hosted edit form.
    const wrapper = await mountPage("presets_Detailing", { id: "1", overlay: "Dialog" });

    // the detailing rendered its host: the form is NOT loaded yet
    const html = wrapper.html();
    expect(html).toContain("data-event=\"edit\"");
    expect(formOpen(wrapper, "Editing Product")).toBe(false);

    // 1) click Edit -> editing.show = true -> the form loads in ONE request
    const before = fetchCount;
    await wrapper.find('[data-event="edit"]').trigger("click");
    await settle();
    expect(fetchCount - before).toBe(1);
    // it is the edit form of record #1
    expect(formOpen(wrapper, "Editing Product")).toBe(true);

    // 2) editing.show = false -> the form is destroyed completely
    await closeForm(wrapper, "Editing Product");
    await settle(200);
    expect(formOpen(wrapper, "Editing Product")).toBe(false);

    // 3) clicking Edit again loads a FRESH form (another single request)
    const before2 = fetchCount;
    await wrapper.find('[data-event="edit"]').trigger("click");
    await settle();
    expect(fetchCount - before2).toBe(1);
    expect(formOpen(wrapper, "Editing Product")).toBe(true);

    wrapper.unmount();
  }, 30000);
});

describe("form host — NEW (listing)", () => {
  it("the New button only flips creating.show, and that opens/destroys the form", async () => {
    const wrapper = await mountPage("presets_OpenListingDialog", {});

    // the listing rendered its host: the create form is NOT loaded yet
    expect(wrapper.html()).toContain("data-event=\"new\"");
    expect(formOpen(wrapper, "New Product")).toBe(false);

    // 1) click New -> creating.show = true -> the form loads in ONE request
    const before = fetchCount;
    await wrapper.find('[data-event="new"]').trigger("click");
    await settle();
    expect(fetchCount - before).toBe(1);
    expect(formOpen(wrapper, "New Product")).toBe(true);

    // 2) creating.show = false -> destroyed
    await closeForm(wrapper, "New Product");
    await settle(200);
    expect(formOpen(wrapper, "New Product")).toBe(false);

    // 3) New again loads a fresh form
    const before2 = fetchCount;
    await wrapper.find('[data-event="new"]').trigger("click");
    await settle();
    expect(fetchCount - before2).toBe(1);
    expect(formOpen(wrapper, "New Product")).toBe(true);

    wrapper.unmount();
  }, 30000);
});
