// DOM integration test (bun + happy-dom) for the EditForm form-scope wrapper:
// mounts the real corejs Root, clicks an Edit button that fires presets_EditForm
// against the Go server (in-memory SQLite), and follows the whole REACTIVE flow
// from the click — the closer turns on, the run-script fires the inner Edit, the
// dialog mounts with the `form` populated from the database.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { nextTick } from "vue";
import { flushPromises } from "@vue/test-utils";
import { installFetch, makeEditButton, mountPresets, type ServerFetch } from "./dom";
import { startServer, type TestServer } from "./server";

let server: TestServer;
let sf: ServerFetch;

beforeAll(async () => {
  server = await startServer("listeditor");
  sf = installFetch(server);
});

afterAll(() => {
  sf?.restore();
  server?.stop();
});

const EditButton = makeEditButton("/admin/products");

describe("EditForm DOM flow", () => {
  it("click -> wrapper -> closer -> inner Edit dialog with fields from DB", async () => {
    const wrapper = mountPresets(`<EditButton></EditButton>`, { EditButton });
    await nextTick();

    // click fires EditForm; the outer portal receives the wrapper
    await wrapper.find("#editBtn").trigger("click");
    await flushPromises();
    await nextTick();

    // the closer setup (~100ms) mounts the form and fires the inner Edit
    await new Promise((r) => setTimeout(r, 800));
    await flushPromises();
    await nextTick();

    const html = wrapper.html();
    // the wrapper landed in the outer portal, and the run-script fired the inner
    // Edit into a nested portal
    expect(html).toContain('id="portal--outerPortal"');
    expect((html.match(/id="portal--_\d+"/g) ?? []).length).toBeGreaterThan(0);
    // the inner Edit rendered the dialog for record #1, bound to the shared closer
    expect(html).toContain('data-title="Editing Product Product#1"');

    // with the vx-dialog stub open, the real fields render from the DB — this is
    // the reactive `form` seeded from the server response
    const values = wrapper
      .findAll('input[type="text"]')
      .map((i: any) => (i.element as HTMLInputElement).value);
    expect(values).toContain("P1"); // Name
    expect(values).toContain("A"); // item 0 Label
    expect(values).toContain("B"); // item 1 Label

    wrapper.unmount();
  }, 20000);
});
