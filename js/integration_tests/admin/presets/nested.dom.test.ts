// DOM integration tests (bun + happy-dom) for a FOUR-LEVEL nested list editor
// (L0 → L1 → L2 → L3), driven from the Edit click against the Go server. They
// mount the real corejs Root (with a vx-dialog stub so the fields render) and
// assert on the reactive `form` — read through the inputs (v-model='form[key]')
// and through the body serialized on submit.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { nextTick } from "vue";
import { flushPromises } from "@vue/test-utils";
import { installFetch, makeEditButton, mountPresets, type ServerFetch } from "./dom";
import { startServer, type TestServer } from "./server";

let server: TestServer;
let sf: ServerFetch;

beforeAll(async () => {
  server = await startServer("nested");
  sf = installFetch(server);
});

afterAll(() => {
  sf?.restore();
  server?.stop();
});

const EditButton = makeEditButton("/admin/l0s");

async function openForm() {
  const wrapper = mountPresets(`<EditButton></EditButton>`, { EditButton });
  await nextTick();
  await wrapper.find("#editBtn").trigger("click");
  await flushPromises();
  await nextTick();
  await new Promise((r) => setTimeout(r, 900));
  await flushPromises();
  await nextTick();
  return wrapper;
}

function textInputs(wrapper: any): string[] {
  return wrapper
    .findAll('input[type="text"]')
    .map((i: any) => (i.element as HTMLInputElement).value);
}

function saveButton(wrapper: any) {
  return wrapper
    .findAll("button")
    .find((b: any) => b.attributes("data-event") === "presets_Update");
}

describe("nested list editor (4 levels) — DOM", () => {
  it("opens the form with every level's `form` value from the DB", async () => {
    const wrapper = await openForm();
    const values = textInputs(wrapper);
    expect(values).toContain("root"); // L0
    expect(values).toContain("a"); // L1
    expect(values).toContain("a1"); // L2
    expect(values).toContain("a1x"); // L3 (deepest)
    wrapper.unmount();
  }, 20000);

  it("serializes the whole nested `form` (every level) on submit", async () => {
    const wrapper = await openForm();
    await saveButton(wrapper)!.trigger("click");
    await flushPromises();

    const form = sf.lastUpdate();
    expect(form).toBeTruthy();
    expect(form!.get("Name")).toBe("root");
    expect(form!.get("L1s[0].Name")).toBe("a");
    expect(form!.get("L1s[0].L2s[0].Name")).toBe("a1");
    expect(form!.get("L1s[0].L2s[0].L3s[0].Name")).toBe("a1x");
    wrapper.unmount();
  }, 20000);

  it("clearing a required Name at the deepest level fails validation and re-renders", async () => {
    const wrapper = await openForm();

    // clear the L3 (deepest) required Name (its value is "a1x")
    const l3 = wrapper
      .findAll('input[type="text"]')
      .find((i: any) => (i.element as HTMLInputElement).value === "a1x");
    expect(l3).toBeTruthy();
    await l3!.setValue("");

    await saveButton(wrapper)!.trigger("click");
    await flushPromises();
    await new Promise((r) => setTimeout(r, 400));
    await flushPromises();
    await nextTick();

    // the cleared value reached the server as part of the reactive `form`
    expect(sf.lastUpdate()!.get("L1s[0].L2s[0].L3s[0].Name")).toBe("");
    // the form re-renders with the required error at the nested level
    expect(wrapper.html()).toContain("This field is required");
    wrapper.unmount();
  }, 20000);
});
