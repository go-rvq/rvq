package js

import "embed"

// Vendor are copies of what the pages of a site take from CDNs — Vue,
// Tailwind's script, Font Awesome, the Inter font —, at the versions they
// ask for: served by the app (js/api/vendor) when a CDN does not answer.
//
//go:embed vendor
var Vendor embed.FS
