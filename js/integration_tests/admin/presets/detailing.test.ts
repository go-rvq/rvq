// HTTP-contract tests for the DETAILING's edit-form host.
//
// Both entry points must host the edit form the same way — the detailing EVENT
// (overlay) and the detailing PAGE (defaultPageFunc): they declare the `editing`
// scope var (the edit overlay's closer), guard the form block with it, and the
// Edit button only turns it on.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { eventFunc, portalBody, startServer, type TestServer } from "./helpers";

const MODEL_URL = "/admin/products";
let server: TestServer;

beforeAll(async () => {
  server = await startServer("listeditor");
});

afterAll(() => {
  server?.stop();
});

// hostAssertions holds what every detailing entry point must render.
function expectEditHost(body: string) {
  // the host owns the edit overlay's closer, closed by default
  expect(body).toContain("$presetsEditing: {show:false}");
  // the form block is guarded by it (turning it off destroys the form)
  expect(body).toContain("v-if='$presetsEditing?.show'");
  // the Edit button only turns it on — it carries no plaid of its own
  expect(body).toContain("@click='$presetsEditing.show = true'");
  // the host loads the form itself, binding the overlay to its closer
  expect(body).toContain('eventFunc("presets_Edit")');
  expect(body).toContain("scope({closer: $presetsEditing})");
  expect(body).toContain('query("presets_closer_provided", "true")');
}

describe("detailing edit-form host", () => {
  it("the detailing EVENT hosts the edit form", async () => {
    const r = await eventFunc(server, "presets_Detailing", {
      url: MODEL_URL,
      query: { id: "1", overlay: "Dialog", target_portal: "p" },
    });
    expectEditHost(portalBody(r));
  });

  it("the detailing PAGE (defaultPageFunc) hosts it the same way", async () => {
    // the page is plain HTML (not an event response); its body is carried inside
    // the root portal's :content attribute, so entities must be decoded first.
    const res = await fetch(`${server.url}${MODEL_URL}/1`);
    expect(res.status).toBe(200);
    const html = (await res.text())
      .replaceAll("&#39;", "'")
      .replaceAll("&quot;", '"')
      .replaceAll("\\u003e", ">")
      .replaceAll('\\"', '"');
    expectEditHost(html);
  });
});
