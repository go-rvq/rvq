// The session dies under an open page, and the page survives it.
//
// Every plaid() request says who it is (the X-Plaid-Request header). A
// middleware that finds the session gone answers such a request with 401 and
// `loginURI` — nothing else — so the page is left untouched. plaid() then asks
// for that URI, which brings the login up, and hands it `onLoginSuccess` in its
// scope. Calling that replays the request the user originally made.
//
// No server here: the fetch is stubbed, so what is under test is exactly the
// corejs side of the contract (see admin/login/dialog.go for the other half).

import { afterEach, beforeEach, describe, expect, it } from "bun:test";
// @ts-expect-error corejs @ alias resolved by corejs tsconfig
import { plaid } from "../../../corejs/src/builder";

type Call = { url: string; header: string | null; method: string };

const w = globalThis as any;
let calls: Call[] = [];
let realFetch: any;
let portalBodies: Record<string, string>;

// answers[i] is the answer to the i-th request
function stubFetch(answers: { status: number; body: any }[]) {
  w.fetch = async (url: string, opts: any) => {
    calls.push({
      url,
      method: opts?.method ?? "GET",
      header: opts?.headers?.["X-Plaid-Request"] ?? null,
    });
    const a = answers[calls.length - 1] ?? { status: 200, body: {} };
    return {
      ok: a.status >= 200 && a.status < 300,
      status: a.status,
      statusText: "",
      redirected: false,
      url,
      json: async () => a.body,
    };
  };
}

async function waitFor(cond: () => boolean, ms = 500) {
  const until = Date.now() + ms;
  while (!cond() && Date.now() < until) {
    await new Promise((r) => setTimeout(r, 5));
  }
}

beforeEach(() => {
  calls = [];
  portalBodies = {};
  realFetch = w.fetch;
  // the login portal, at the layout root
  w.__goplaid = {
    portals: {
      presets_LoginPortalName: {
        updatePortalTemplate: (body: string) => {
          portalBodies["presets_LoginPortalName"] = body;
        },
      },
    },
  };
  window.history.replaceState(null, "", "/admin/products");
});

afterEach(() => {
  w.fetch = realFetch;
});

// what the user was doing when the session ended
function interruptedRequest() {
  return plaid()
    .vars({})
    .url("/admin/products")
    .eventFunc("presets_Update")
    .form({ Name: "meio digitado" });
}

describe("session lost under an open page", () => {
  it("opens the login and replays the interrupted request", async () => {
    stubFetch([
      { status: 401, body: { loginURI: "/admin/login-dialog" } },
      // the dialog: its script is what the login calls once the session is back
      {
        status: 200,
        body: {
          updatePortals: [{ name: "presets_LoginPortalName", body: "<div>login</div>" }],
          runScript: "onLoginSuccess()",
        },
      },
      { status: 200, body: {} },
    ]);

    await interruptedRequest().go();
    await waitFor(() => calls.length === 3);

    expect(calls.length).toBe(3);

    // 1) the request that hit the dead session
    expect(calls[0].url).toContain("/admin/products");
    expect(calls[0].url).toContain("__execute_event__=presets_Update");

    // 2) the login, at the address the 401 gave
    expect(calls[1].url).toContain("/admin/login-dialog");
    expect(portalBodies["presets_LoginPortalName"]).toBe("<div>login</div>");

    // 3) the user's action goes through, unchanged
    expect(calls[2].url).toBe(calls[0].url);

    // and the server always knows it is the page asking
    expect(calls.map((c) => c.header)).toEqual(["1", "1", "1"]);
    expect(calls.map((c) => c.method)).toEqual(["POST", "POST", "POST"]);
  });

  it("acts on nothing else in the 401 response", async () => {
    stubFetch([
      {
        status: 401,
        body: {
          loginURI: "/admin/login-dialog",
          // a middleware has no business updating the page — and if it tries,
          // the page is not touched
          updatePortals: [{ name: "presets_LoginPortalName", body: "<div>nope</div>" }],
          redirectURL: "/admin/login",
          pageTitle: "login",
        },
      },
      { status: 200, body: {} },
    ]);

    const title = document.title;
    await interruptedRequest().go();
    await waitFor(() => calls.length === 2);

    expect(portalBodies["presets_LoginPortalName"]).toBeUndefined();
    expect(document.title).toBe(title);
    expect(window.location.pathname).toBe("/admin/products");
  });

  // A 401 WITHOUT a login URI keeps falling through to the usual error message —
  // that path is not covered here because it rejects a promise nobody awaits
  // (it always has), which aborts the run under `bun test`.
});
