// Prism not highlighting the page by itself (before its core runs): the
// blocks of a text are highlighted by mdkit's code (code.js), and only them.
if (typeof window !== 'undefined') {
  window.Prism = window.Prism || {}
  window.Prism.manual = true
}
