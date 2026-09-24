// HTTP contract for the POST-SAVE REFRESH (the DOM behaviour is in
// refresh.dom.test.ts).
//
// A saved form does not know who opened it, and carries nothing about it in its
// request. It ends by calling every hook in `onSaveCallbacks` — a list that
// travels down the scope like `closer` and `form`, each form host appending its
// own before passing it on:
//
//   listing hosts       → reload the listing
//   detailing (overlay) → render itself again in the portal it occupies; the
//                         listing that opened it is already in the list
//   detailing (page)    → ReloadDetail: re-renders its portal and returns the
//                         new page title, so nothing is reloaded

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { eventFunc, portalBody, startServer, updateSigned, type TestServer } from "./helpers";

let server: TestServer;

beforeAll(async () => {
  server = await startServer("refresh");
});

afterAll(() => {
  server?.stop();
});

const unesc = (s: string) => s.replace(/\\u003e/g, ">").replace(/\\u0026/g, "&");

async function page(path: string) {
  return unesc(await (await fetch(`${server.url}${path}`)).text());
}

describe("saving hands over to the closer", () => {
  it("an update ends by calling onSave with the saved id, and needs no callback query", async () => {
    const r = await updateSigned(server, { url: "/admin/articles", query: { id: "1", overlay: "Dialog" }, fields: { Title: "A1 edited", Body: "b1" } });

    expect(r.runScript ?? "").toContain('(onSaveCallbacks || []).forEach(f => f("1"))');
    // the request carried no presets_post_change_callback at all
    expect(JSON.stringify(r)).not.toContain("presets_post_change_callback");
  });

  it("a create ends the same way, with the NEW record's id", async () => {
    const r = await eventFunc(server, "presets_Create", {
      method: "POST",
      url: "/admin/articles",
      query: { overlay: "Dialog" },
      fields: { Title: "created", Body: "x" },
    });

    const script = r.runScript ?? "";
    expect(script).toMatch(/onSaveCallbacks \|\| \[\]\)\.forEach\(f => f\("\d+"\)\)/);
    expect(script).not.toContain('f("")');
  });
});

describe("what each host refreshes", () => {
  it("a listing's hosts (new, item edit, item detail) reload the listing", async () => {
    const body = await page("/admin/articles");

    // the three hosts of a listing, each APPENDING its own refresh to the list
    for (const scope of ["$presetsCreating", "$presetsItemEditing", "$presetsItemDetailing"]) {
      expect(body).toContain(scope);
    }
    expect(body.split("[...onSaveCallbacks, (id) =>").length - 1).toBeGreaterThanOrEqual(3);
    expect(body).toContain("reload: () =>");
    expect(body).toContain("presets_ReloadList");
  });

  it("a listing without detailing still reloads on the row's edit", async () => {
    const body = await page("/admin/tags");

    // no detail host, but the edit one refreshes the listing all the same
    expect(body).not.toContain("$presetsItemDetailing");
    expect(body).toContain("$presetsItemEditing");
    expect(body).toContain("presets_ReloadList");
  });

  it("an overlayed detailing refreshes itself and passes the news up", async () => {
    const r = await eventFunc(server, "presets_Detailing", {
      url: "/admin/articles",
      query: { id: "1", overlay: "Dialog", target_portal: "p" },
    });
    const body = unesc(portalBody(r, "p"));

    // it appends "render me again" to the list — not delegated to whoever opened
    // it, which may not be a form host at all; the listing that opened it is
    // already in the list with its own hook
    expect(body).toContain("[...onSaveCallbacks, (id) =>");
    expect(body).toContain("presets_Detailing");
  });

  it("a detailing PAGE refreshes its portal instead, never the page", async () => {
    const body = await page("/admin/articles/1");

    expect(body).toContain("presets_ReloadDetail");
    expect(body).toContain("presets_detail_page");
    expect(body).toContain("[...onSaveCallbacks, (id) =>");
  });
});

describe("ReloadDetail", () => {
  it("re-renders the detail body with the new data and the new title", async () => {
    await updateSigned(server, { url: "/admin/articles", query: { id: "2", overlay: "Dialog" }, fields: { Title: "A2 renamed", Body: "b2" } });

    const r = await eventFunc(server, "presets_ReloadDetail", {
      url: "/admin/articles",
      query: { id: "2" },
    });

    // the page title follows the record, and the client puts it in <title>
    expect(r.pageTitle).toBe("A2 renamed");
    expect(portalBody(r, "presets_detail_page")).toContain("A2 renamed");

    // and nothing asks the browser to reload
    expect(r.reload).toBeFalsy();
    expect(r.pushState).toBeFalsy();
  });

  it("works for a SINGLETON detailing (no id)", async () => {
    await updateSigned(server, { url: "/admin/site", query: { id: "1", overlay: "Dialog" }, fields: { Name: "Site renamed", Motto: "m" } });

    const r = await eventFunc(server, "presets_ReloadDetail", { url: "/admin/site" });
    expect(portalBody(r, "presets_detail_page")).toContain("Site renamed");
    expect(r.reload).toBeFalsy();
  });
});
