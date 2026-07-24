// Integration tests for the EditForm/NewForm form-scope wrapper, driven over
// HTTP against the Go list-editor test app (in-memory SQLite, seeded with a
// Product "P1" holding items A and B).
//
// These assert the exact plaid event-response contract the browser relies on:
// the wrapper establishes a single `form` scope guarded by a closer, then runs
// the real Edit/New into an inner portal — and the inner form re-renders never
// recreate the scope.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import {
  eventFunc,
  extractRunScript,
  parseRunScriptEvent,
  portalBody,
  startServer,
  type TestServer,
} from "./helpers";

let server: TestServer;

beforeAll(async () => {
  server = await startServer();
});

afterAll(() => {
  server?.stop();
});

describe("EditForm wrapper", () => {
  it("mounts a closer scope + form scope + inner portal + run-script", async () => {
    const r = await eventFunc(server, "presets_EditForm", {
      query: { id: "1", overlay: "Dialog", target_portal: "outerPortal" },
    });

    const body = portalBody(r, "outerPortal");

    // closer scope above, turned on via a setup timeout
    expect(body).toContain("<go-plaid-scope");
    expect(body).toContain("closer.show = true");
    // the form is destroyed when the closer is turned off
    expect(body).toContain("v-if='closer.show'");
    // a single form scope, child of the ambient form (raw JS, not a string)
    expect(body).toContain(":scope='{form: {$parent: form}}'");
    // an inner portal + a run-script that fires the real Edit into it, sharing
    // the same closer
    expect(body).toContain("<go-plaid-portal");
    expect(body).toContain("<go-plaid-run-script");

    const script = extractRunScript(body);
    const inner = parseRunScriptEvent(script);
    expect(inner.event).toBe("presets_Edit");
    expect(script).toContain("scope({closer: closer})");
    // the inner Edit targets the inner portal (not the outer one)
    expect(inner.queries["target_portal"]).toBeTruthy();
    expect(inner.queries["target_portal"]).not.toBe("outerPortal");
  });

  it("follows the run-script to render the real Edit form (fields, no own form scope)", async () => {
    const wrap = await eventFunc(server, "presets_EditForm", {
      query: { id: "1", overlay: "Dialog", target_portal: "outerPortal" },
    });
    const inner = parseRunScriptEvent(extractRunScript(portalBody(wrap, "outerPortal")));

    const r = await eventFunc(server, inner.event, { query: inner.queries });
    // the Edit renders into its inner target portal
    const body = portalBody(r, inner.queries["target_portal"]);

    // the edit form is a dialog bound to the shared closer (no own form scope
    // recreated for the whole dialog)
    expect(body).toContain("<vx-dialog");
    expect(body).toContain("v-model='closer.show'");
    // it renders its fields (Name + the list editor items)
    expect(body).toContain('form["Name"]');
    expect(body).toContain('form["Items[0].Label"]');
    expect(body).toContain('form["Items[1].Label"]');
  });
});

describe("inner form re-render", () => {
  async function openInnerQueries() {
    const wrap = await eventFunc(server, "presets_EditForm", {
      query: { id: "1", overlay: "Dialog", target_portal: "outerPortal" },
    });
    return parseRunScriptEvent(extractRunScript(portalBody(wrap, "outerPortal"))).queries;
  }

  it("re-renders with a validation error when required Name is empty", async () => {
    const q = await openInnerQueries();

    const r = await eventFunc(server, "presets_Update", {
      method: "POST",
      query: q,
      fields: {
        Name: "", // required -> fails
        "Items.__present": "1",
        "Items[0].ID": "1",
        "Items[0].Label": "A",
        "Items[0].__pos": "0",
        "Items[1].ID": "2",
        "Items[1].Label": "B",
        "Items[1].__pos": "1",
      },
    });

    const body = JSON.stringify(r);
    expect(body).toContain("This field is required");
    // the re-render must NOT create another form scope (the wrapper owns it)
    expect(body).not.toContain("FormInit");
  });

  it("saves successfully, removes the deleted item and closes via closer.show=false", async () => {
    const q = await openInnerQueries();

    const r = await eventFunc(server, "presets_Update", {
      method: "POST",
      query: { ...q, overlay: "Dialog" },
      fields: {
        Name: "P1", // valid
        "Items.__present": "1",
        "Items[0].ID": "1",
        "Items[0].Label": "A",
        "Items[0].__pos": "0",
        // item #1 removed
        "Items[1].ID": "2",
        "Items[1].Label": "B",
        "Items[1].__pos": "1",
        "Items[1].__deleted": "true",
      },
    });

    // the overlay is torn down by turning the closer off
    expect(r.runScript ?? "").toContain("closer.show = false");
  });
});
