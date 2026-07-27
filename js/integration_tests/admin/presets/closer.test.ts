// The closer contract: `show` fires open/close callbacks on TRANSITION, and a
// closer that carries a clean page `url` puts it in the address bar while it is
// open — restoring the previous one when it closes, LIFO.
//
// No server here: this is the corejs unit that every overlay is built on.

import { beforeEach, describe, expect, it } from "bun:test";
// @ts-expect-error corejs @ alias resolved by corejs tsconfig
import { asCloser, createCloser, isCloser } from "../../../corejs/src/closer";
// @ts-expect-error corejs @ alias resolved by corejs tsconfig
import { closerUrlStack, resetCloserUrlStack } from "../../../corejs/src/closer-url";

beforeEach(() => {
  resetCloserUrlStack();
  window.history.replaceState(null, "", "/admin/articles");
});

describe("open/close callbacks", () => {
  it("fires only on transition, and hands the closer over", () => {
    const seen: string[] = [];
    const c = createCloser({});
    c.onOpen((closer: any) => seen.push("open:" + (closer === c)));
    c.onClose((closer: any) => seen.push("close:" + (closer === c)));

    c.show = true;
    c.show = true; // same value: nothing happens
    c.show = false;
    c.show = false;

    expect(seen).toEqual(["open:true", "close:true"]);
  });

  it("does not fire for the state the initializer carried", () => {
    let opens = 0;
    const c = createCloser({ show: true });
    c.onOpen(() => opens++);

    expect(c.show).toBe(true);
    expect(opens).toBe(0);

    c.show = false;
    c.show = true;
    expect(opens).toBe(1);
  });

  it("keeps going when a callback throws", () => {
    const c = createCloser({});
    let reached = false;
    c.openCallbacks.push(() => {
      throw new Error("boom");
    });
    c.onOpen(() => {
      reached = true;
    });

    c.show = true;
    expect(reached).toBe(true);
  });

  it("onOpen/onClose return the remover", () => {
    let opens = 0;
    const c = createCloser({});
    const off = c.onOpen(() => opens++);

    c.show = true;
    off();
    c.show = false;
    c.show = true;

    expect(opens).toBe(1);
  });

  it("asCloser is idempotent — adopting one must not duplicate its callbacks", () => {
    const c = createCloser({});
    let opens = 0;
    c.onOpen(() => opens++);

    expect(isCloser(c)).toBe(true);
    expect(asCloser(c)).toBe(c);
    expect(c.openCallbacks.length).toBe(2); // the url sync + ours

    c.show = true;
    expect(opens).toBe(1);
  });
});

describe("address bar, LIFO", () => {
  it("shows the overlay's clean URL and restores the previous one", () => {
    const detail = createCloser({ url: "/admin/articles/1" });

    detail.show = true;
    expect(window.location.pathname).toBe("/admin/articles/1");

    detail.show = false;
    expect(window.location.pathname).toBe("/admin/articles");
    expect(closerUrlStack().length).toBe(0);
  });

  it("stacks and unwinds in order — the deepest closes first", () => {
    const detail = createCloser({ url: "/admin/articles/1" });
    const edit = createCloser({ url: "/admin/articles/1/edit" });

    detail.show = true;
    edit.show = true;
    expect(window.location.pathname).toBe("/admin/articles/1/edit");
    expect(closerUrlStack().map((e: any) => e.url)).toEqual([
      "/admin/articles/1",
      "/admin/articles/1/edit"
    ]);

    edit.show = false;
    expect(window.location.pathname).toBe("/admin/articles/1");

    detail.show = false;
    expect(window.location.pathname).toBe("/admin/articles");
  });

  it("closing an overlay that has others above it takes them along", () => {
    const listing = createCloser({ url: "/admin/articles" });
    const detail = createCloser({ url: "/admin/articles/1" });

    window.history.replaceState(null, "", "/admin/home");
    listing.show = true;
    detail.show = true;
    expect(window.location.pathname).toBe("/admin/articles/1");

    // the outer one goes away (its dialog is destroyed, and with it the inner)
    listing.show = false;
    expect(window.location.pathname).toBe("/admin/home");
    expect(closerUrlStack().length).toBe(0);
  });

  it("takes the address from a function of the closer — one host, many rows", () => {
    const host = createCloser({
      id: null,
      url: (c: any) => "/admin/articles/" + c.id
    });

    host.id = "7";
    host.show = true;
    expect(window.location.pathname).toBe("/admin/articles/7");

    host.show = false;
    host.id = "9";
    host.show = true;
    expect(window.location.pathname).toBe("/admin/articles/9");
    host.show = false;
  });

  it("a closer with no url leaves the address bar alone", () => {
    const c = createCloser({});
    c.show = true;
    expect(window.location.pathname).toBe("/admin/articles");
    expect(closerUrlStack().length).toBe(0);
    c.show = false;
  });

  it("Back closes the top overlay instead of leaving a stale address", () => {
    const detail = createCloser({ url: "/admin/articles/1" });
    const edit = createCloser({ url: "/admin/articles/1/edit" });

    detail.show = true;
    edit.show = true;

    // the browser moved on its own: our handler closes what that address showed
    window.dispatchEvent(new Event("popstate"));
    expect(edit.show).toBe(false);
    expect(detail.show).toBe(true);
    expect(closerUrlStack().map((e: any) => e.url)).toEqual(["/admin/articles/1"]);

    window.dispatchEvent(new Event("popstate"));
    expect(detail.show).toBe(false);
    expect(closerUrlStack().length).toBe(0);
  });
});
