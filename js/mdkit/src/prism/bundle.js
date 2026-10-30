// The Prism of a page that shows a text (built to rvq's
// js/vendor/prism: loaded by mdkit's code.js when a block of code is there
// and the page has no Prism of its own): the languages and plugins of the
// admin's code viewer (vuetifyx vx-code) — line numbers, a toolbar with the
// language and a copy button — (their look is code.js's). window.Prism.
import './manual.js'
import Prism from 'prismjs'
import 'prismjs/components/prism-clike'
import 'prismjs/components/prism-markup'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-markdown'
import 'prismjs/components/prism-ini'
import 'prismjs/components/prism-toml'
import 'prismjs/components/prism-sql'
import 'prismjs/components/prism-diff'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-docker'
import 'prismjs/components/prism-makefile'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-jsx'
import 'prismjs/components/prism-tsx'
import 'prismjs/plugins/line-numbers/prism-line-numbers'
import 'prismjs/plugins/toolbar/prism-toolbar'
import 'prismjs/plugins/show-language/prism-show-language'
import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard'
import { registerGad, registerGadx, registerGadTemplate } from '@gad-lang/prism-gad'

// registerGad first: gadx and gadt embed it
registerGad(Prism)
registerGadx(Prism)
registerGadTemplate(Prism)
Prism.languages.golang = Prism.languages.go
Prism.languages.gadtemplate = Prism.languages.gadt

window.Prism = Prism
