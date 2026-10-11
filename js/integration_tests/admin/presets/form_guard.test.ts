// The guard of a form in an overlay: a close the user asks for
// (closer.requestClose) is held — confirmClose on — while the CONTENT of the
// form differs from what each field held first; a close by the app
// (`show = false`, once saved) is not guarded.
//
// No server here: the corejs units (closer, form-guard).

import { describe, expect, it } from "bun:test";
// @ts-expect-error corejs @ alias resolved by corejs tsconfig
import { createCloser } from "../../../corejs/src/closer";
// @ts-expect-error corejs @ alias resolved by corejs tsconfig
import { FormBaseline, fieldsOf, keyOf } from "../../../corejs/src/form-guard";

describe("closer guards", () => {
  it("closes when no guard holds", () => {
    const c = createCloser({ show: true });
    expect(c.requestClose()).toBe(true);
    expect(c.show).toBe(false);
    expect(c.confirmClose).toBe(false);
  });

  it("holds the close a guard says, and asks", () => {
    const c = createCloser({ show: true });
    let dirty = true;
    const remove = c.onCloseGuard(() => dirty);

    expect(c.requestClose()).toBe(false);
    expect(c.show).toBe(true);
    expect(c.confirmClose).toBe(true);

    // discarded: the app closes; the confirmation is over
    c.show = false;
    expect(c.confirmClose).toBe(false);

    c.show = true;
    dirty = false;
    expect(c.requestClose()).toBe(true);
    remove();
    expect(c.closeGuards.length).toBe(0);
  });

  it("does not guard a close by the app (show = false: saved)", () => {
    const c = createCloser({ show: true });
    c.onCloseGuard(() => true);
    c.show = false;
    expect(c.show).toBe(false);
  });

  it("goes on when a guard throws, holding nothing for it", () => {
    const c = createCloser({ show: true });
    c.onCloseGuard(() => {
      throw new Error("boom");
    });
    expect(c.requestClose()).toBe(true);
  });
});

describe("the content of the form", () => {
  it("compares each field with what it held first", () => {
    const form: any = { $parent: {}, Name: "a", Count: 1 };
    const b = new FormBaseline();
    b.see(form);
    expect(b.changed(form)).toBe(false);

    form.Name = "b";
    expect(b.changed(form)).toBe(true);
    form.Name = "a";
    expect(b.changed(form)).toBe(false);
  });

  it("takes a field that comes later (a tab drawn on open) as its first value", () => {
    const form: any = { Name: "a" };
    const b = new FormBaseline();
    b.see(form);
    form.Late = "from the server";
    b.see(form);
    expect(b.changed(form)).toBe(false);
    form.Late = "typed";
    expect(b.changed(form)).toBe(true);
  });

  it("leaves the scope's keys out", () => {
    expect(fieldsOf({ $parent: 1, __x: 2, Name: 3 })).toEqual(["Name"]);
  });

  it("tells files apart by name, size and date", () => {
    const f1 = new File(["ab"], "a.txt", { lastModified: 1 });
    const f2 = new File(["abc"], "a.txt", { lastModified: 1 });
    expect(keyOf(f1)).not.toBe(keyOf(f2));
  });
});
