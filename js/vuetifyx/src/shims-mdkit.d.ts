// rvq mdkit (js/mdkit: '@rvq/mdkit', its source) is JavaScript: its exports,
// untyped.
declare module '@rvq/mdkit' {
  export const FootnoteBuilder: any
  export const TableBuilder: any
  export const RichEditor: any
  export const MarkdownEditor: any
  export const FootnoteRef: any
  export const Footnotes: any
  export const finalizeFootnotes: (html: string) => string
  export const tiptapFootnoteAdapter: (editor: any, refPos?: number | null) => any
  export const tables: any
  export const footnotes: any
  export const inlineHTML: (md: string) => string
  export function configure(options: { textClass?: string, words?: Record<string, Record<string, string>> }): any
  export function textStyle(cls?: string): void
  export function footnoteBalloons(): void
  export function install(app: any, options?: { prefix?: string }): void
  export const components: Record<string, any>
  export function linkCards(root: Element, resolve: (targets: string[]) => Promise<Record<string, any>>): Promise<void>
  export function watchLinkCards(root?: Element, resolve?: ((targets: string[]) => Promise<Record<string, any>>) | null): void
  export function metaResolver(doc?: Document): ((targets: string[]) => Promise<Record<string, any>>) | null
}
