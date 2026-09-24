// HTTP helpers for the presets integration tests: drive the Go test server's
// event funcs and assert on the plaid event responses. Server boot lives in
// ./server (works under both bun and vitest).

import { type TestServer } from "./server";

export { startServer, type TestServer } from "./server";

export interface PortalUpdate {
  name: string;
  body: string;
}

export interface EventResponse {
  body?: string;
  runScript?: string;
  updatePortals?: PortalUpdate[];
  pushState?: unknown;
  reload?: boolean;
}

const EDIT_URL = "/admin/products";

// eventFunc issues a plaid event-func request and returns the decoded response.
// GET events carry their params in the query; POST events also send `fields` as
// the (url-encoded) form body. url-encoding — not FormData — keeps the request
// runtime-agnostic (the DOM tests' happy-dom preload replaces the global FormData
// with one that drops empty fields); the Go server accepts non-multipart submits.
export async function eventFunc(
  server: TestServer,
  event: string,
  opts: {
    method?: "GET" | "POST";
    query?: Record<string, string>;
    fields?: Record<string, string>;
    url?: string;
  } = {},
): Promise<EventResponse> {
  const method = opts.method ?? "GET";
  const u = new URL((opts.url ?? EDIT_URL), server.url);
  u.searchParams.set("__execute_event__", event);
  for (const [k, v] of Object.entries(opts.query ?? {})) {
    u.searchParams.set(k, v);
  }

  const init: RequestInit = { method };
  if (opts.fields) {
    const usp = new URLSearchParams();
    for (const [k, v] of Object.entries(opts.fields)) usp.append(k, v);
    init.body = usp;
  }

  const res = await fetch(u.toString(), init);
  const text = await res.text();
  if (!res.headers.get("content-type")?.includes("application/json")) {
    throw new Error(`event ${event}: non-JSON response (${res.status}):\n${text}`);
  }
  return JSON.parse(text) as EventResponse;
}

// recordStamp is the signed record stamp an edit form seeds (`__formSign`), read
// from the form the way the browser has it. An update without it is refused as
// a stale form ("This form is out of date"): the stamp is what proves the form
// was rendered over the record as it stands.
export async function recordStamp(
  server: TestServer,
  opts: { url?: string; id?: string },
): Promise<Record<string, string>> {
  const r = await eventFunc(server, "presets_Edit", {
    method: "POST",
    url: opts.url,
    query: { id: opts.id ?? "", overlay: "Dialog" },
  });
  const html = (r.updatePortals ?? []).map((p) => p.body).join("") + (r.body ?? "");
  const sign = formAssigns(html)["__formSign"];
  return sign ? { __formSign: String(sign) } : {};
}

// updateSigned posts a presets_Update as the browser does: with the record
// stamp of a freshly opened form beside the fields.
export async function updateSigned(
  server: TestServer,
  opts: {
    method?: "GET" | "POST";
    query?: Record<string, string>;
    fields?: Record<string, string>;
    url?: string;
  },
): Promise<EventResponse> {
  const stamp = await recordStamp(server, { url: opts.url, id: opts.query?.id });
  return eventFunc(server, "presets_Update", {
    ...opts,
    method: "POST",
    fields: { ...stamp, ...(opts.fields ?? {}) },
  });
}

// portalBody returns the body of the named updated portal (or the whole response
// body when name is omitted / not found).
export function portalBody(r: EventResponse, name?: string): string {
  if (name) {
    const p = r.updatePortals?.find((p) => p.name === name);
    if (p) return p.body;
  }
  return r.updatePortals?.[0]?.body ?? r.body ?? "";
}

// extractRunScript returns the script body of the (first) <go-plaid-run-script>
// element in an HTML fragment (the wrapper emits it as :script='(scope) => {…}').
export function extractRunScript(html: string): string {
  const m = html.match(/<go-plaid-run-script :script='([\s\S]*?)'>/);
  if (!m) throw new Error("no <go-plaid-run-script> in:\n" + html);
  return m[1];
}

// followEditForm fires the EditForm/NewForm wrapper and follows its run-script to
// render the real inner Edit/New, returning the inner event response. This mirrors
// what the browser does on the Edit click.
export async function followEditForm(
  server: TestServer,
  url: string,
  id: string,
): Promise<EventResponse> {
  const wrap = await eventFunc(server, "presets_EditForm", {
    url,
    query: { id, overlay: "Dialog", target_portal: "outer" },
  });
  const inner = parseRunScriptEvent(extractRunScript(portalBody(wrap, "outer")));
  return eventFunc(server, inner.event, { url, query: inner.queries });
}

// formAssigns extracts the reactive form's initial values from the rendered
// fragment: each input seeds form["<key>"] via v-assign='[form, {"key": val}]',
// so the merged object is what the reactive `form` holds.
export function formAssigns(html: string): Record<string, any> {
  const out: Record<string, any> = {};
  const re = /v-assign='\[form, (\{.*?\})\]'/g;
  let m: RegExpExecArray | null;
  while ((m = re.exec(html))) {
    try {
      Object.assign(out, JSON.parse(m[1]));
    } catch {
      /* skip non-JSON */
    }
  }
  return out;
}

// parseRunScriptEvent extracts the inner event func + queries from a run-script
// (the `plaid()...eventFunc("x").queries({...}).go()` the EditForm/NewForm wrapper
// emits), so a test can follow it and drive the inner form.
export function parseRunScriptEvent(script: string): {
  event: string;
  queries: Record<string, string>;
} {
  const ev = script.match(/eventFunc\("([^"]+)"\)/);
  if (!ev) throw new Error("no eventFunc() in script:\n" + script);
  const q: Record<string, string> = {};
  const qm = script.match(/queries\((\{.*?\})\)/);
  if (qm) {
    const obj = JSON.parse(qm[1]) as Record<string, string[]>;
    for (const [k, v] of Object.entries(obj)) q[k] = Array.isArray(v) ? v[0] : String(v);
  }
  return { event: ev[1], queries: q };
}
