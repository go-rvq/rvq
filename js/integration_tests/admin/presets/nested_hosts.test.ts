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
  it("a page's hosts live on `vars`; a dialog's never do", async () => {
    const page = await pageHTML("/admin/products");
    const detailPage = await pageHTML("/admin/products/1");
    const dialog = portalBody(
      await eventFunc(server, "presets_OpenListingDialog", { query: { target_portal: "p" } }),
      "p",
    );

    for (const body of [page, detailPage, dialog]) {
      expect(body).toContain("$presets");
    }

    // On a PAGE the state goes to `vars`: that is what gives it an address any
    // portal can reach, and without one the overlay would have to be rendered
    // inside the host's own portal — where a drawer sizes itself against
    // whatever box holds it and opens below the app bar.
    expect(page).toContain("vars.$presetsCreating");
    expect(page).toContain("vars.$presetsItemDetailing");
    expect(page).toContain("vars.$presetsItemEditing");

    // A dialog keeps slot variables: two live listings sharing `vars` is what
    // made them fight over the same name, mounting every guarded block at once.
    expect(dialog).not.toContain("vars.$presets");
  });

  it("a listing in a dialog declares its OWN state, like the page does", async () => {
    const dialog = portalBody(
      await eventFunc(server, "presets_OpenListingDialog", { query: { target_portal: "p" } }),
      "p",
    );

    // its own scope + guarded blocks, addressed as plain slot variables
    expect(dialog).toContain("$presetsItemDetailing: $closer({show:false, id:null,");
    expect(dialog).toContain("v-if='$presetsItemDetailing?.show'");
    expect(dialog).toContain('query("id", $presetsItemDetailing.id)');
  });

  it("a detailing opened as an overlay hosts its edit form on its own state", async () => {
    const r = await eventFunc(server, "presets_Detailing", {
      query: { id: "1", overlay: "Dialog", target_portal: "p" },
    });
    const body = portalBody(r, "p");

    expect(body).toContain("$presetsEditing: $closer({show:false,");
    expect(body).toContain("v-if='$presetsEditing?.show'");
    expect(body).toContain("@click='$presetsEditing.show = true'");
    expect(body).not.toContain("vars.$presetsEditing");
  });
});
