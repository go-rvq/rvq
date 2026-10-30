# vendor

Copies of what the pages of a site take from CDNs, at fixed versions: the
fallback when a CDN does not answer (`js/api/vendor`: each asset's tag asks
the CDN first and, on error, these, served at `/_rvq/vendor/`).

| Asset | Version | From | License |
|---|---|---|---|
| `vue/3.5.44/vue.global.prod.js` | 3.5.44 | unpkg.com/vue | MIT |
| `tailwind/3.4.17/tailwind.js` | 3.4.17 | cdn.tailwindcss.com (Play CDN) | MIT |
| `fontawesome/6.5.1/` (css, webfonts) | 6.5.1 Free | cdnjs | icons CC BY 4.0, fonts SIL OFL 1.1, code MIT |
| `inter/` (css, woff2) | Google Fonts, weights 300–800 | fonts.googleapis.com | SIL OFL 1.1 |
| `katex/0.16.47/` (js, css, fonts woff2) | 0.16.47 | cdn.jsdelivr.net/npm/katex | MIT |
| `mermaid/11.17.2/mermaid.min.js` | 11.17.2 | cdn.jsdelivr.net/npm/mermaid | MIT |
| `prism/1.30.0/prism.min.js` | Prism 1.30.0, its languages and plugins, prism-gad 0.1.0-rc.2 | built by `js/mdkit/build.sh` (no CDN: served as `?v=` its hash) | MIT |

A new version: download it beside the old one, point the asset at it
(`js/api/vendor`), and its integrity with it.
