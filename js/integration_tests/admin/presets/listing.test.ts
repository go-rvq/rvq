// HTTP-contract tests for the LISTING's form hosts: the listing declares ONE
// host per action and the rows only flip it, instead of every row carrying its
// own overlay plaid.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { eventFunc, portalBody, startServer, type TestServer } from "./helpers";

const URL = "/admin/products";
let server: TestServer;

beforeAll(async () => {
  server = await startServer("listeditor");
});

afterAll(() => {
  server?.stop();
});

describe("listing form hosts", () => {
  it("declares the shared per-item and create hosts", async () => {
    const r = await eventFunc(server, "presets_OpenListingDialog", {
      url: URL,
      query: { target_portal: "p" },
    });
    const body = portalBody(r);

    // one host per action, each owning its overlay's closer
    expect(body).toContain("$presetsItemDetailing: {show:false, id:null,");
    expect(body).toContain("$presetsCreating: {show:false,");
    // each host guards its block, so turning the scope off destroys the overlay
    expect(body).toContain("v-if='$presetsItemDetailing?.show'");
    expect(body).toContain("v-if='$presetsCreating?.show'");
    // and loads through the shared portal, keyed by the scope's id
    expect(body).toContain('query("id", $presetsItemDetailing.id)');
    expect(body).toContain('query("presets_closer_provided", "true")');
  });

  it("rows open the shared host for their record", async () => {
    const r = await eventFunc(server, "presets_OpenListingDialog", {
      url: URL,
      query: { target_portal: "p" },
    });
    const body = portalBody(r);

    // the row click only points the host at its record and turns it on …
    expect(body).toContain('$presetsItemDetailing.id = "1"; $presetsItemDetailing.show = true');
    // … carrying no plaid of its own (the request is issued by the host's
    // run-script, once, when the scope is turned on)
    const rowClick = body.match(/@click\.self='([^']*)'/)?.[1] ?? "";
    expect(rowClick).toContain("$presetsItemDetailing.show = true");
    expect(rowClick).not.toContain("plaid()");
  });

  it("the New button only flips the create host", async () => {
    const r = await eventFunc(server, "presets_OpenListingDialog", {
      url: URL,
      query: { target_portal: "p" },
    });
    const body = portalBody(r);

    expect(body).toContain("@click='$presetsCreating.show = true'");
    expect(body).not.toContain('eventFunc("presets_NewForm")');
  });
});
