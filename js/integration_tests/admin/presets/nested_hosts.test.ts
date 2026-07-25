// HTTP contract for NESTED form hosts (the DOM behaviour is in
// nested_hosts.dom.test.ts).
//
// Every render declares its host state as a SLOT variable, never on the
// app-global `vars` — that is what lets several levels be alive at once (page
// listing, a nested list in a dialog, a detail opened from it) without watching
// each other's variable.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { eventFunc, portalBody, startServer, type TestServer } from "./helpers";

let server: TestServer;

beforeAll(async () => {
  server = await startServer();
});

afterAll(() => {
  server?.stop();
});

async function pageHTML(path: string) {
  const r = await fetch(`${server.url}${path}`);
  return await r.text();
}

describe("nested form hosts — contract", () => {
  it("no host state is published on `vars` (page or overlay)", async () => {
    const page = await pageHTML("/admin/products");
    const detailPage = await pageHTML("/admin/products/1");
    const dialog = portalBody(
      await eventFunc(server, "presets_OpenListingDialog", { query: { target_portal: "p" } }),
      "p",
    );

    for (const body of [page, detailPage, dialog]) {
      expect(body).toContain("$presets");
      // sharing `vars` is exactly what made two live listings fight over the
      // same variable, mounting every guarded block at once
      expect(body).not.toContain("vars.$presets");
    }
  });

  it("a listing in a dialog declares its OWN state, like the page does", async () => {
    const dialog = portalBody(
      await eventFunc(server, "presets_OpenListingDialog", { query: { target_portal: "p" } }),
      "p",
    );

    // its own scope + guarded blocks, addressed as plain slot variables
    expect(dialog).toContain("$presetsItemDetailing: {show:false, id:null}");
    expect(dialog).toContain("v-if='$presetsItemDetailing?.show'");
    expect(dialog).toContain('query("id", $presetsItemDetailing.id)');
  });

  it("a detailing opened as an overlay hosts its edit form on its own state", async () => {
    const r = await eventFunc(server, "presets_Detailing", {
      query: { id: "1", overlay: "Dialog", target_portal: "p" },
    });
    const body = portalBody(r, "p");

    expect(body).toContain("$presetsEditing: {show:false}");
    expect(body).toContain("v-if='$presetsEditing?.show'");
    expect(body).toContain("@click='$presetsEditing.show = true'");
    expect(body).not.toContain("vars.$presetsEditing");
  });
});
