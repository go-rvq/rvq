// rvq-portal: a place of the page whose content is a Vue template given at
// run time (README.md). The library: the component, its types, the
// compilation of a template into a component (component-by-template); and
// install, which registers it as <rvq-portal>.
import type { App } from 'vue'

import Portal from './Portal'

export { Portal }
export * from './types'
export * from './component-by-template'

export function install(app: App) {
  app.component('rvq-portal', Portal)
}
