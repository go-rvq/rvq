// A LISTAGEM EM PÁGINA — o caso que os outros testes de DOM não cobriam, porque
// todos abrem a listagem dentro de um diálogo.
//
// Numa página, o botão "+" é renderizado pelo LAYOUT na barra de ações, e o
// overlay pedido é o RightDrawer (e não o Dialog). O corpo da página viaja
// escapado no portal do app, então o teste monta exatamente o que o navegador
// monta.

import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { nextTick } from "vue";
import { flushPromises } from "@vue/test-utils";
import { installFetch, mountPresets, type ServerFetch } from "./dom";
import { startServer, type TestServer } from "./server";

let server: TestServer;
let sf: ServerFetch;
const requests: string[] = [];

beforeAll(async () => {
  server = await startServer("refresh");
  sf = installFetch(server);

  const real = globalThis.fetch;
  globalThis.fetch = ((input: any, init?: any) => {
    const url = typeof input === "string" ? input : input?.url;
    if (typeof url === "string") requests.push(url);
    return real(input, init);
  }) as typeof fetch;
});

afterAll(() => {
  sf?.restore();
  server?.stop();
});

async function settle(ms = 700) {
  await flushPromises();
  await nextTick();
  await new Promise((r) => setTimeout(r, ms));
  await flushPromises();
  await nextTick();
}

// O corpo de uma página é o payload do portal do app: o parser de HTML decodifica
// as entidades do atributo, e o que sobra é a string (JSON) com o template.
function pageTemplate(html: string): string {
  const m = html.match(/raw :content='("(?:\\.|[^"\\])*")'/);
  if (!m) {
    throw new Error("portal do app não encontrado na página");
  }
  const decoded = m[1]
    .replace(/&#39;/g, "'")
    .replace(/&quot;/g, '"')
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&amp;/g, "&");
  return JSON.parse(decoded);
}

async function mountPage(path: string) {
  const html = await (await fetch(server.url + path)).text();
  const wrapper = mountPresets(pageTemplate(html));
  await settle();
  return wrapper;
}

describe("NEW a partir de uma listagem em página", () => {
  it("o + da barra de ações abre o formulário no drawer", async () => {
    const wrapper = await mountPage("/admin/articles");

    const btn = wrapper.find('[data-event="new"]');
    expect(btn.exists()).toBe(true);

    // o layout já tem o seu (o menu), então o que conta é o do FORMULÁRIO
    const formDrawer = () =>
      wrapper.findAll(".vx-drawer").find((d: any) => d.find("form").exists() || d.text().includes("Title"));
    expect(formDrawer()).toBeUndefined();

    requests.length = 0;
    await btn.trigger("click");
    await settle();

    // um pedido, e do formulário de criação
    const events = requests.filter((r) => r.includes("__execute_event__"));
    expect(events.length).toBe(1);
    expect(events[0]).toContain("presets_New");
    expect(events[0]).toContain("overlay=RightDrawer");

    // e o drawer ABRIU: a resposta tem de cair no portal do host, cujo escopo
    // tem o closer que o botão ligou. No portal global do layout o `closer` é o
    // da raiz, e o drawer ficaria fechado para sempre.
    const opened = formDrawer();
    expect(opened).toBeDefined();
    expect(opened!.text()).toContain("Title");

    wrapper.unmount();
  }, 40000);

  it("a linha da tabela abre o detalhe no drawer", async () => {
    const wrapper = await mountPage("/admin/articles");

    const formDrawer = () =>
      wrapper.findAll(".vx-drawer").find((d: any) => d.text().includes("Title"));
    expect(formDrawer()).toBeUndefined();

    requests.length = 0;
    await wrapper.findAll("table")[0].findAll("td")[0].trigger("click");
    await settle();

    const events = requests.filter((r) => r.includes("__execute_event__"));
    expect(events.length).toBe(1);
    expect(events[0]).toContain("presets_Detailing");
    expect(formDrawer()).toBeDefined();

    wrapper.unmount();
  }, 40000);
});
