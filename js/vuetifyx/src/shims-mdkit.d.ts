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
  export const LinkCard: any
  export function watchBlocks(root?: Element, assets?: Record<string, any>): void
  export function drawBlocks(root: Element, assets?: Record<string, any>): Promise<any>
  export function loadVendor(name: string, global: string, assets?: Record<string, any>): Promise<any>
  export const Diagram: any
  export const InlineMath: any
  export const BlockMath: any
  export const MathCommands: any
  export function drawDiagramInto(el: Element, source: string): Promise<boolean>
  export const LinkCardPicker: any
  export const HelpDialog: any
  export const LinkPicker: any
  export const LinkForm: any
  export function hrefOf(target: string): Promise<string>
  export function cardTarget(text: string, origin?: string): string
  export function metaSearch(doc?: Document): ((q: string) => Promise<any[]>) | null
  export function linkCards(root: Element, resolve: (targets: string[]) => Promise<Record<string, any>>): Promise<void>
  export function watchLinkCards(root?: Element, resolve?: ((targets: string[]) => Promise<Record<string, any>>) | null): void
  export function metaResolver(doc?: Document): ((targets: string[]) => Promise<Record<string, any>>) | null
}
