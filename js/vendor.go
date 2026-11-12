package js

import "embed"

// Vendor are copies of what the pages of a site take from CDNs — Vue,
// Tailwind's script, Font Awesome, the Inter font —, at the versions they
// ask for: served by the app (js/api/vendor) when a CDN does not answer.
//
// The folder is cdn, not vendor: a module's zip leaves out what is under a
// subfolder of a folder named vendor (the vendored packages of old) — a
// dependent of rvq got none of them.
//
//go:embed cdn
var Vendor embed.FS
