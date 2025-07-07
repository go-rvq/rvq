import {ckeditor} from "@/lib/ckeditor.ts";

declare const window: any;
window.__goplaidVueComponentRegisters =
  window.__goplaidVueComponentRegisters || [];
window.__goplaidVueComponentRegisters.push(ckeditor.install);
