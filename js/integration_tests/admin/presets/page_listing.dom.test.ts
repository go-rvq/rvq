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

// O REFRESH depois do save, numa listagem em página. Aqui os overlays vão para o
// portal do LAYOUT e alcançam o estado do host por referência
// (`vars.$presets…`). Os ganchos de refresh (`onSaveCallbacks`) têm de ir pelo
// mesmo caminho: no portal do layout, o `onSaveCallbacks` do escopo é a lista
// vazia da raiz, e um save ali não recarregaria nada.

function tableText(wrapper: any) {
  return wrapper.findAll("table").map((t: any) => t.text()).join(" | ");
}

// Digita no primeiro campo do formulário aberto e salva. O formulário é o do
// botão de salvar mais recente — o drawer que o contém —, e um input oculto não
// é campo (o primeiro é a assinatura do registro, `__formSign`).
async function fillAndSave(wrapper: any, value: string) {
  const save = [
    ...wrapper.findAll('[data-event="presets_Update"]'),
    ...wrapper.findAll('[data-event="presets_Create"]'),
  ];
  expect(save.length).toBeGreaterThan(0);
  const btn = save[save.length - 1];

  const drawer = (btn.element as HTMLElement).closest(".vx-drawer, .vx-dialog, .v-navigation-drawer, .v-dialog") as HTMLElement;
  expect(drawer).toBeTruthy();
  const input = drawer.querySelector('input:not([type="hidden"])') as HTMLInputElement;
  expect(input).toBeTruthy();
  input.value = value;
  input.dispatchEvent(new Event("input", { bubbles: true }));
  await settle(200);

  requests.length = 0;
  await btn.trigger("click");
  await settle();
}

describe("refresh depois do save, numa listagem em página", () => {
  it("NEW: a listagem mostra o registro criado", async () => {
    const wrapper = await mountPage("/admin/articles");
    await wrapper.find('[data-event="new"]').trigger("click");
    await settle();

    await fillAndSave(wrapper, "Page created");

    expect(requests.some((r) => r.includes("presets_ReloadList"))).toBe(true);
    expect(tableText(wrapper)).toContain("Page created");
    wrapper.unmount();
  }, 40000);

  it("EDIT da linha (modelo sem detalhe): a listagem mostra o valor novo", async () => {
    const wrapper = await mountPage("/admin/tags");
    await wrapper.findAll("table")[0].findAll("td")[0].trigger("click");
    await settle();

    await fillAndSave(wrapper, "Page tag edited");

    expect(requests.some((r) => r.includes("presets_ReloadList"))).toBe(true);
    expect(tableText(wrapper)).toContain("Page tag edited");
    wrapper.unmount();
  }, 40000);

  it("DETAIL → EDIT: o detalhe que abriu o form E a listagem recarregam", async () => {
    const wrapper = await mountPage("/admin/articles");
    await wrapper.findAll("table")[0].findAll("td")[0].trigger("click");
    await settle();

    await wrapper.find('[data-event="edit"]').trigger("click");
    await settle();
    await fillAndSave(wrapper, "Page detail edited");

    // o detalhe se redesenha no portal que ocupa…
    expect(requests.some((r) => r.includes("presets_ReloadDetail") || r.includes("presets_Detailing"))).toBe(true);
    const detail = wrapper.findAll(".vx-drawer").find((d: any) => d.text().includes("Page detail edited") && !d.find("table").exists());
    expect(detail).toBeDefined();
    // …e a listagem embaixo também
    expect(tableText(wrapper)).toContain("Page detail edited");
    wrapper.unmount();
  }, 40000);
});
