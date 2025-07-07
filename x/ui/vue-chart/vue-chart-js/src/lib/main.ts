import { vuechart } from "@/lib/vuechart.ts";

declare const window: any;
window.__goplaidVueComponentRegisters =
  window.__goplaidVueComponentRegisters || [];
window.__goplaidVueComponentRegisters.push(vuechart.install);
