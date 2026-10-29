// The plugins of vx-tiptap-editor: by name (its prop `plugins`, the Go
// VXTipTapEditorBuilder.Plugins and TiptapBuilder.Plugins), each its
// extensions — with their buttons on the bar — and what it makes of the HTML
// the field keeps (output). Registered here: footnotes, tables (rvq mdkit's
// builders), cards (its link cards); an application registers its own (registerTiptapPlugin).
import type {AnyExtension} from '@tiptap/core'

export interface TiptapPluginContext {
  /** t is the words of the namespace vx (useI18n) */
  t: (key: string) => string
}

export interface TiptapPlugin {
  /** extensions are the plugin's extensions of the editor */
  extensions?: (ctx: TiptapPluginContext) => AnyExtension[]
  /** output is the HTML the field keeps, from the editor's */
  output?: (html: string) => string
}

const registry = new Map<string, TiptapPlugin>()

/** registerTiptapPlugin registers (or replaces) the plugin name. */
export function registerTiptapPlugin(name: string, plugin: TiptapPlugin) {
  registry.set(name, plugin)
}

/** tiptapPlugins are the plugins of names registered, in their order. */
export function tiptapPlugins(names: readonly string[] = []): TiptapPlugin[] {
  return names.map((n) => registry.get(n)).filter((p): p is TiptapPlugin => !!p)
}

// the plugins of rvq: the footnotes and the tables of mdkit
import footnotes from './footnotes'
import tablesPlugin from './tables'
import cards from './cards'

registerTiptapPlugin('footnotes', footnotes)
registerTiptapPlugin('tables', tablesPlugin)
registerTiptapPlugin('cards', cards)
