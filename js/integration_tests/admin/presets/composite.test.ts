// Integration tests for a multi-level nested list editor over records with
// COMPOSITE primary keys (the association-object flavor of many-to-many):
//
//   Cart (ID)
//    └─ Items[]   CartItem   PK (CartID, Sku)          — required Qty
//        └─ Notes[] CartNote PK (CartID, Sku, Seq)     — required Text
//
// Driven over HTTP against the Go server (in-memory SQLite). They follow the Edit
// click, assert the reactive `form` (read from the v-assign seeds), and exercise
// validation, inclusion, deletion and re-render at composite-key nesting levels —
// keying each row by its FULL composite primary key.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import {
  eventFunc,
  followEditForm,
  formAssigns,
  portalBody,
  startServer,
  type TestServer,
} from "./helpers";

const URL = "/admin/carts";
let server: TestServer;

beforeAll(async () => {
  server = await startServer("composite");
});

afterAll(() => {
  server?.stop();
});

// the seeded item and note, as full field sets (all composite PK columns), so a
// test can post a valid form and tweak one thing.
function seededFields(): Record<string, string> {
  return {
    Name: "c1",
    "Items.__present": "1",
    "Items[0].CartID": "1",
    "Items[0].Sku": "A",
    "Items[0].Qty": "2",
    "Items[0].__pos": "0",
    "Items[0].__index": "0",
    "Items[0].Notes.__present": "1",
    "Items[0].Notes[0].CartID": "1",
    "Items[0].Notes[0].Sku": "A",
    "Items[0].Notes[0].Seq": "1",
    "Items[0].Notes[0].Text": "hello",
    "Items[0].Notes[0].__pos": "0",
    "Items[0].Notes[0].__index": "0",
  };
}

describe("composite-key nested list editor", () => {
  it("opens with every composite-key level's `form` from the DB", async () => {
    const r = await followEditForm(server, URL, "1");
    const form = formAssigns(portalBody(r));

    expect(form["Name"]).toBe("c1");
    // level 1: CartItem, PK (CartID, Sku). Hidden uint keys seed as numbers; the
    // editable Qty is a text field, so it seeds as a string.
    expect(form["Items[0].CartID"]).toBe(1);
    expect(form["Items[0].Sku"]).toBe("A");
    expect(form["Items[0].Qty"]).toBe("2");
    // level 2: CartNote, PK (CartID, Sku, Seq)
    expect(form["Items[0].Notes[0].CartID"]).toBe(1);
    expect(form["Items[0].Notes[0].Sku"]).toBe("A");
    expect(form["Items[0].Notes[0].Seq"]).toBe(1);
    expect(form["Items[0].Notes[0].Text"]).toBe("hello");
  });

  it("clearing a required field at the deepest composite level fails validation", async () => {
    const fields = seededFields();
    fields["Items[0].Notes[0].Text"] = ""; // required Text at the deepest level
    const r = await eventFunc(server, "presets_Update", {
      method: "POST",
      url: URL,
      query: { id: "1", overlay: "Dialog", target_portal: "inner" },
      fields,
    });

    expect(JSON.stringify(r)).toContain("This field is required");
    // the re-render preserves the composite-key values of the other fields
    const form = formAssigns(portalBody(r));
    expect(form["Items[0].Sku"]).toBe("A");
    expect(form["Items[0].Qty"]).toBe("2");
    expect(form["Items[0].Notes[0].Seq"]).toBe(1);
  });

  it("adds a note (inclusion) at a composite-key level and flags it __new", async () => {
    const r = await eventFunc(server, "listEditor_addRowEvent", {
      method: "POST",
      url: URL,
      query: { id: "1", listEditor_AddRowFormKey: "Items[0].Notes" },
      fields: seededFields(),
    });

    const body = portalBody(r);
    // the appended note (index 1) is flagged __new by the Setup re-seed
    expect(body).toContain('const newKeys = ["Items[0].Notes[1]"]');
    // the existing note is still present
    expect(body).toContain("Items[0].Notes[0].Text");
  });

  it("removing a deep composite-key note keeps it deleted on a validation re-render", async () => {
    const fields = seededFields();
    fields["Name"] = ""; // fail validation to force a re-render
    fields["Items[0].Notes[0].__deleted"] = "true"; // remove the deepest note
    const r = await eventFunc(server, "presets_Update", {
      method: "POST",
      url: URL,
      query: { id: "1", overlay: "Dialog", target_portal: "inner" },
      fields,
    });

    const body = portalBody(r);
    expect(JSON.stringify(r)).toContain("This field is required"); // root Name
    // the removed deep note is re-seeded as deleted (keyed by its composite key row)
    expect(body).toContain('const deletedKeys = ["Items[0].Notes[0]"]');
  });

  it("removing a first-level composite-key item persists on a valid save", async () => {
    // the root reconciles its Items has-many (composite PK CartID,Sku), so
    // removing the only item persists.
    const r = await eventFunc(server, "presets_Update", {
      method: "POST",
      url: URL,
      query: { id: "1", overlay: "Dialog", target_portal: "inner" },
      fields: {
        Name: "c1",
        "Items.__present": "1",
        "Items[0].CartID": "1",
        "Items[0].Sku": "A",
        "Items[0].Qty": "2",
        "Items[0].__pos": "0",
        "Items[0].__deleted": "true", // remove the whole composite-key item
      },
    });

    // save succeeded -> the dialog closes
    expect(r.runScript ?? "").toContain("closer.show = false");

    // reopening shows no items anymore
    const reopened = await followEditForm(server, URL, "1");
    const form = formAssigns(portalBody(reopened));
    expect(form["Name"]).toBe("c1");
    expect(form["Items[0].Sku"]).toBeUndefined();
  });
});
