/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import {registerPlugins} from "@/demo/plugins";
import {vuechart} from "@/lib/vuechart";

// Components
import App from "@/demo/App.vue";

// Composables
import {createApp} from "vue";

const app = createApp(App);

registerPlugins(app);
app.use(vuechart);
app.mount("#app");
