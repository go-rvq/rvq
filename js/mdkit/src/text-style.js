import { config } from './config.js'

// textStyle is the look of a text written in the editors, once for the page
// and the class cls (config textClass) — the editors' own, and where the text
// is shown —: headings of every level, lists, quotes, code, rules, tables,
// footnotes (and their balloon), sections that fold; a page with no typography of its own (a
// site's Tailwind) has it so.
export function textStyle(cls = config.textClass) {
  const id = 'mdkit-text-style-' + cls
  if (document.getElementById(id)) return
  const st = document.createElement('style')
  st.id = id
  st.textContent = [
    '.' + cls + '>*+*{margin-top:.75em}',
    '.' + cls + ' h1{font-size:1.6em;font-weight:700;line-height:1.25;margin-top:1em}',
    '.' + cls + ' h2{font-size:1.35em;font-weight:700;line-height:1.3;margin-top:1em}',
    '.' + cls + ' h3{font-size:1.15em;font-weight:600;line-height:1.3;margin-top:.9em}',
    '.' + cls + ' h4{font-size:1.05em;font-weight:600;line-height:1.35;margin-top:.9em}',
    '.' + cls + ' h5{font-size:1em;font-weight:600;line-height:1.4;margin-top:.8em}',
    '.' + cls + ' h6{font-size:.9em;font-weight:600;line-height:1.4;margin-top:.8em;color:#4b5563}',
    '.' + cls + ' hr{border:0;border-top:1px solid #d1d5db;margin:1.25em 0}',
    // a table fluid (the default): as wide as the text (its widths parts or
    // least widths of it); not (mdkit-fit): as wide as its columns, placed
    // left, center or right; on a narrow screen, scrolled
    '.' + cls + ' table{border-collapse:collapse}',
    '.' + cls + ' table:not(.mdkit-fit){width:100% !important}',
    '.' + cls + ' table.mdkit-fit{width:auto !important}',
    '.' + cls + ' table.mdkit-center{margin-left:auto;margin-right:auto}',
    '.' + cls + ' table.mdkit-right{margin-left:auto;margin-right:0}',
    '@media (max-width:640px){.' + cls + ' table{display:block;overflow-x:auto}}',
    '.' + cls + ' th,.' + cls + ' td{border:1px solid #d1d5db;padding:.35em .75em}',
    '.' + cls + ' th{background:#f3f4f6;font-weight:600}',
    '.' + cls + '>:first-child{margin-top:0}',
    '.' + cls + ' ul{list-style:disc;padding-left:1.5em}',
    '.' + cls + ' ol{list-style:decimal;padding-left:1.5em}',
    '.' + cls + ' li+li{margin-top:.25em}',
    '.' + cls + ' blockquote{border-left:3px solid #d1d5db;padding-left:1em;color:#4b5563}',
    '.' + cls + ' pre{background:#f3f4f6;border-radius:.375rem;padding:.75em 1em;overflow-x:auto;font-size:.875em}',
    '.' + cls + ' code{font-family:ui-monospace,monospace;font-size:.9em}',
    '.' + cls + ' :not(pre)>code{background:#f3f4f6;border-radius:.25rem;padding:.1em .3em}',
    '.' + cls + ' a{color:#2563eb;text-decoration:underline}',
    '.' + cls + ' img{max-width:100%;height:auto}',
    // footnotes: the reference small and up, the notes at the end apart, the
    // way back quiet, the note reached by a link marked a moment
    '.' + cls + ' sup{font-size:.75em;line-height:0}',
    '.' + cls + ' a.footnote-ref{text-decoration:none;font-weight:600;padding:0 .1em}',
    '.' + cls + ' .footnotes{font-size:.875em;color:#4b5563;margin-top:1.5em}',
    '.' + cls + ' .footnotes hr{border:0;border-top:1px solid #d1d5db;margin:0 0 .75em;width:30%}',
    '.' + cls + ' .footnotes ol{list-style:decimal;padding-left:1.5em}',
    '.' + cls + ' .footnotes li p{margin:0}',
    '.' + cls + ' a.footnote-backref{text-decoration:none;margin-left:.25em}',
    '.' + cls + ' li:target,.' + cls + ' sup:target{background:#fef9c3;transition:background 1s}',
    // a section that folds: its title a button — a chevron that turns, a
    // hover —; its content in, a line on its left
    '.' + cls + ' details{margin:.75em 0}',
    '.' + cls + ' details>summary{cursor:pointer;list-style:none;display:flex;align-items:center;gap:.5em;font-weight:600;' +
      'padding:.3em .5em;margin-left:-.5em;border-radius:6px;width:fit-content;max-width:100%}',
    '.' + cls + ' details>summary::-webkit-details-marker{display:none}',
    '.' + cls + ' details>summary::before{content:"";flex:0 0 auto;width:.45em;height:.45em;border-right:2px solid currentColor;' +
      'border-bottom:2px solid currentColor;transform:rotate(-45deg);transition:transform .15s}',
    '.' + cls + ' details[open]>summary::before{transform:rotate(45deg)}',
    '.' + cls + ' details>summary:hover{background:rgba(128,128,128,.12)}',
    '.' + cls + ' details>:not(summary){margin-left:.35em;padding-left:1em;border-left:3px solid #d1d5db}',
    '.' + cls + ' details>:not(summary)+:not(summary){margin-top:0;padding-top:.75em}',
    // open, a card (as a tab printed): its title the header, its content in
    '.' + cls + ' details[open]{border:1px solid rgba(128,128,128,.35);border-radius:6px;padding:0 1em .75em}',
    '.' + cls + ' details[open]>summary{width:auto;max-width:none;margin:0 -1em .75em;padding:.5em 1em;border-radius:6px 6px 0 0;' +
      'background:rgba(128,128,128,.06);border-bottom:1px solid rgba(128,128,128,.25)}',
    '.' + cls + ' details[open]>:not(summary){margin-left:0;padding-left:0;border-left:0}',
    '.' + cls + ' details[open]>:not(summary)+:not(summary){padding-top:0}',
    // the balloon of a note
    '.mdkit-fn-balloon{position:absolute;z-index:3000;max-width:22rem;background:#fff;color:#111827;border:1px solid #d1d5db;' +
      'border-radius:6px;box-shadow:0 6px 18px rgba(0,0,0,.15);padding:.5rem .75rem;font-size:.875rem;line-height:1.45}',
    '.mdkit-fn-balloon p{margin:0}',
    '.mdkit-fn-balloon a.footnote-backref{display:none}',
  ].join('')
  document.head.append(st)
}
