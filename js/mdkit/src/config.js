// The configuration of mdkit, for the page (configure): the class of a text
// as it is shown (its look: textStyle) — the editors' and their previews' —,
// and the words of the editors when no component is given its own (by
// editor: rich, markdown, table, footnote; else window.mdkitWords, or
// window.hermonWords: a site of hermon's).
export const config = { textClass: 'mdkit-text', words: null }

// configure sets what options says (textClass, words).
export function configure(options = {}) {
  if (options.textClass) config.textClass = options.textClass
  if (options.words) config.words = options.words
  return config
}

// wordsOf is the words of the editor name (rich, markdown, table, footnote).
export function wordsOf(name) {
  const all = config.words || (typeof window !== 'undefined' && (window.mdkitWords || window.hermonWords)) || {}
  return all[name] || {}
}
