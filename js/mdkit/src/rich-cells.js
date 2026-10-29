// inlineHTML is a line of Markdown (a cell of a table) as HTML: its bold,
// italic, code, links, … (marked's inline). What is kept is the server's to
// sanitize, as any text.
import { marked } from 'marked'

export const inlineHTML = (md) => marked.parseInline(String(md || ''))
