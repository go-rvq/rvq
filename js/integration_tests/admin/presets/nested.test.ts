// Integration tests for a FOUR-LEVEL nested list editor (L0 → L1 → L2 → L3),
// driven over HTTP against the Go server (in-memory SQLite). They follow the Edit
// click (EditForm → inner Edit), assert the reactive `form` (read from the
// v-assign seeds the browser binds each input to), and exercise validation,
// inclusion and deletion at nested depth — always checking the submitted `form`.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { eventFunc, followEditForm, formAssigns, portalBody, startServer, type TestServer, updateSigned } from "./helpers";

const URL = "/admin/l0s";
let server: TestServer;

beforeAll(async () => {
  server = await startServer("nested");
});

afterAll(() => {
  server?.stop();
});

describe("nested list editor (4 levels)", () => {
  it("opens the form with the whole nested `form` from the DB", async () => {
    const r = await followEditForm(server, URL, "1");
    const body = portalBody(r);
    const form = formAssigns(body);

    // every nesting level's Name is seeded into the reactive form
    expect(form["Name"]).toBe("root");
    expect(form["L1s[0].Name"]).toBe("a");
    expect(form["L1s[0].L2s[0].Name"]).toBe("a1");
    expect(form["L1s[0].L2s[0].L3s[0].Name"]).toBe("a1x");
  });

  it("clearing a required Name at the deepest level fails validation", async () => {
    const r = await updateSigned(server, {
      method: "POST",
      url: URL,
      query: { id: "1", overlay: "Dialog", target_portal: "inner" },
      fields: {
        Name: "root",
        "L1s.__present": "1",
        "L1s[0].ID": "1",
        "L1s[0].Name": "a",
        "L1s[0].__pos": "0",
        "L1s[0].L2s.__present": "1",
        "L1s[0].L2s[0].ID": "1",
        "L1s[0].L2s[0].Name": "a1",
        "L1s[0].L2s[0].__pos": "0",
        "L1s[0].L2s[0].L3s.__present": "1",
        "L1s[0].L2s[0].L3s[0].ID": "1",
        "L1s[0].L2s[0].L3s[0].Name": "", // required at the deepest level -> fails
        "L1s[0].L2s[0].L3s[0].__pos": "0",
      },
    });

    const body = JSON.stringify(r);
    expect(body).toContain("This field is required");
    // the re-render preserves the other levels' `form` values
    const form = formAssigns(portalBody(r));
    expect(form["Name"]).toBe("root");
    expect(form["L1s[0].Name"]).toBe("a");
    expect(form["L1s[0].L2s[0].Name"]).toBe("a1");
  });

  it("adds a row at a nested level (inclusion) and flags it __new", async () => {
    // add a row to the L2s of L1[0] — the add-row event re-renders that item
    const r = await eventFunc(server, "listEditor_addRowEvent", {
      method: "POST",
      url: URL,
      query: {
        id: "1",
        [`listEditor_AddRowFormKey`]: "L1s[0].L2s",
      },
      fields: {
        Name: "root",
        "L1s.__present": "1",
        "L1s[0].ID": "1",
        "L1s[0].Name": "a",
        "L1s[0].__pos": "0",
        "L1s[0].__index": "0",
        "L1s[0].L2s.__present": "1",
        "L1s[0].L2s[0].ID": "1",
        "L1s[0].L2s[0].Name": "a1",
        "L1s[0].L2s[0].__pos": "0",
        "L1s[0].L2s[0].__index": "0",
      },
    });

    const body = portalBody(r);
    // the appended L2 row (index 1) is flagged __new by the Setup re-seed
    expect(body).toContain('const newKeys = ["L1s[0].L2s[1]"]');
    // the existing L2 is still there
    expect(body).toContain("L1s[0].L2s[0].Name");
  });

  it("removing a nested item keeps it deleted on a validation re-render", async () => {
    // remove the deepest item (L3[0]) AND clear the root Name so validation
    // fails: the re-render must keep the removed L3 flagged deleted (its removal
    // is not lost across the round-trip).
    const r = await updateSigned(server, {
      method: "POST",
      url: URL,
      query: { id: "1", overlay: "Dialog", target_portal: "inner" },
      fields: {
        Name: "", // required -> validation fails, forcing a re-render
        "L1s.__present": "1",
        "L1s[0].ID": "1",
        "L1s[0].Name": "a",
        "L1s[0].__pos": "0",
        "L1s[0].L2s.__present": "1",
        "L1s[0].L2s[0].ID": "1",
        "L1s[0].L2s[0].Name": "a1",
        "L1s[0].L2s[0].__pos": "0",
        "L1s[0].L2s[0].L3s.__present": "1",
        "L1s[0].L2s[0].L3s[0].ID": "1",
        "L1s[0].L2s[0].L3s[0].Name": "a1x",
        "L1s[0].L2s[0].L3s[0].__pos": "0",
        "L1s[0].L2s[0].L3s[0].__deleted": "true", // removed at the deepest level
      },
    });

    const body = portalBody(r);
    // the required error rendered (root Name)
    expect(JSON.stringify(r)).toContain("This field is required");
    // and the removed deepest item is re-seeded as deleted by the Setup
    expect(body).toContain('const deletedKeys = ["L1s[0].L2s[0].L3s[0]"]');
  });

  it("removing a first-level item removes it from the DB on save", async () => {
    // the root reconciles its L1s has-many, so removing the only L1 persists.
    const r = await updateSigned(server, {
      method: "POST",
      url: URL,
      query: { id: "1", overlay: "Dialog", target_portal: "inner" },
      fields: {
        Name: "root",
        "L1s.__present": "1",
        "L1s[0].ID": "1",
        "L1s[0].Name": "a",
        "L1s[0].__pos": "0",
        "L1s[0].__deleted": "true", // remove the whole L1 subtree
      },
    });

    // save succeeded -> the dialog closes
    expect(r.runScript ?? "").toContain("closer.show = false");

    // reopening the form shows no L1 items anymore
    const reopened = await followEditForm(server, URL, "1");
    const form = formAssigns(portalBody(reopened));
    expect(form["Name"]).toBe("root");
    expect(form["L1s[0].Name"]).toBeUndefined();
  });
});
