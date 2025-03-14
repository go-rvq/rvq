import {type App} from "vue";

import {default as CKEditor5} from "./Ckeditor5.vue";

export const ckeditor = {
  install: (app: App) => {
    app.component("ckeditor5", CKEditor5);
  },
};
